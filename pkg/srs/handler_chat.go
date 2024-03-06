package srs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dariallab/srs/pkg/templates"
	"github.com/gorilla/websocket"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

type wsChatInput struct {
	Message string `json:"message"`
}

type wsChatResponse struct {
	Original string
	Diff     []dmp.Diff
	Response string
}

func (s *Server) showChatHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := getUserFromContext(ctx)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t := templates.TemplateChat
	if err := t.Execute(w, nil); err != nil {
		s.logger.Error().Err(err).Msg("can't execute template")
		http.Error(w, errInternal, http.StatusInternalServerError)
		return
	}
}

func (s *Server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := getUserFromContext(ctx)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		return
	}

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

		input := strings.ReplaceAll(in.Message, "\n", "")

		originalOut := &wsChatResponse{
			Original: input,
		}

		var originalTpl bytes.Buffer
		if err := templates.TemplateChatInput.Execute(&originalTpl, nil); err != nil {
			s.logger.Error().Err(err).Msg("can't execute chat input template")
			return
		}
		if err := templates.TemplateChatMessages.Execute(&originalTpl, originalOut); err != nil {
			s.logger.Error().Err(err).Msg("can't execute chat response template with original")
			return
		}

		if err = ws.WriteMessage(websocket.TextMessage, originalTpl.Bytes()); err != nil {
			s.logger.Error().Err(err).Msg("can't write message to web socket")
			return
		}

		corrected, err := s.ai.Correct(ctx, input)
		if err != nil {
			s.logger.Error().Err(err).Msg("can't correct message")
		}

		if corrected != "" {
			diff := Diff(input, corrected)
			out := &wsChatResponse{
				Diff: diff,
			}

			var responseTpl bytes.Buffer
			if err := templates.TemplateChatMessages.Execute(&responseTpl, out); err != nil {
				s.logger.Error().Err(err).Msg("can't execute chat response template with diff")
				return
			}

			if err = ws.WriteMessage(websocket.TextMessage, responseTpl.Bytes()); err != nil {
				s.logger.Error().Err(err).Msg("can't write message to web socket")
				return
			}
		}

		response, err := s.ai.Response(ctx, input)
		if err != nil {
			s.logger.Error().Err(err).Msg("can't respond to the message")
			return
		}

		if response != "" {
			out := &wsChatResponse{
				Response: response,
			}

			var responseTpl bytes.Buffer
			if err := templates.TemplateChatMessages.Execute(&responseTpl, out); err != nil {
				s.logger.Error().Err(err).Msg("can't execute chat response template with response")
				return
			}

			if err = ws.WriteMessage(websocket.TextMessage, responseTpl.Bytes()); err != nil {
				s.logger.Error().Err(err).Msg("can't write message to web socket")
				return
			}
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

func Diff(before, after string) []dmp.Diff {
	d := dmp.New()
	diffs := d.DiffMain(before, after, true)
	return d.DiffCleanupSemantic(diffs)
}
