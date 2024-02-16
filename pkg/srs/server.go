package srs

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/dariallab/srs/pkg/templates"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type wsChatInput struct {
	Message string `json:"message"`
}

type wsChatResponse struct {
	Original  string `json:"original"`
	Corrected string `json:"corrected"`
}

type AI interface {
	Correct(ctx context.Context, input string) (string, error)
}

type Server struct {
	http.Handler
	ai     AI
	logger zerolog.Logger
}

func NewServer(ai AI, logger zerolog.Logger) *Server {
	s := &Server{
		ai:     ai,
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

	t := templates.TemplateChat
	if err := t.Execute(w, nil); err != nil {
		s.logger.Error().Err(err).Msg("can't execute template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

		var in wsChatInput
		if err = json.Unmarshal(m, &in); err != nil {
			s.logger.Error().Err(err).Msg("can't unmarshal message from web socket")
			continue
		}

		corrected, err := s.ai.Correct(ctx, in.Message)
		if err != nil {
			s.logger.Error().Err(err).Msg("can't correct message")
		}

		out := &wsChatResponse{
			Original:  in.Message,
			Corrected: corrected,
		}

		var tpl bytes.Buffer
		if err := templates.TemplateChatResponse.Execute(&tpl, out); err != nil {
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
