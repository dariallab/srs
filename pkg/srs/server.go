package srs

import (
	"bytes"
	"encoding/json"
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
	ws, err := upgradeToWebSocket(w, r)
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
			continue
		}

		var tpl bytes.Buffer
		if err := static.TemplateChatResponse.Execute(&tpl, in); err != nil {
			s.logger.Error().Err(err).Msg("can't execute chat response template")
			continue
		}

		if err = ws.WriteMessage(websocket.TextMessage, tpl.Bytes()); err != nil {
			s.logger.Error().Err(err).Msg("can't write message to web socket")
			return
		}
	}
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	return ws, err
}
