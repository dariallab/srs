package srs

import (
	"encoding/json"
	"fmt"
	"net/http"

	static "github.com/dariallab/srs/pkg/templates"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type wsMessage struct {
	Message string `json:"message"`
}

type Server struct {
	http.Handler
	logger zerolog.Logger
}

func NewServer(logger zerolog.Logger) *Server {
	s := &Server{
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /chat", s.showChatHandler)
	mux.HandleFunc("GET /message", s.sendMessageHandler)
	s.Handler = mux

	return s
}

func (s *Server) showChatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t := static.TemplateChat
	if err := t.Execute(w, nil); err != nil {
		s.logger.Error().Err(err).Msg("can't execute template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't upgrade to web socket")
		return
	}
	defer ws.Close()

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error().Err(err).Msg("unexpected close of web socket")
			}
			return
		}

		var in wsMessage
		if err = json.Unmarshal(m, &in); err != nil {
			s.logger.Error().Err(err).Msg("can't unmarshal message from web socket")
			return
		}

		if err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`
		<div id="chat-response" hx-swap-oob="beforeend"><p>%s</p></div>
		<input id="input" type="text" name="message" placeholder="Type your message here" required autofocus>
		`, in.Message))); err != nil {
			s.logger.Error().Err(err).Msg("can't write message to web socket")
			return
		}

	}
}
