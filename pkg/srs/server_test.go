package srs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestTutorServer(t *testing.T) {
	server := NewServer(zerolog.New(os.Stdout))

	t.Run("open page - see form", func(t *testing.T) {
		req := newRequest(t, http.MethodGet, "/chat", nil)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), `<form ws-send>`, "websocket tags are not found")
	})

	t.Run("send message - return response via websockets", func(t *testing.T) {
		server := httptest.NewServer(server)
		defer server.Close()
		msg := `{"message":"hello"}`

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/message"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("can't open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			t.Fatalf("can't send message over ws connection %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("can't read the websocket message: %s ", err)
		}

		assert.Contains(t, string(p), `<div id="chat-response" hx-swap-oob="beforeend"><p>hello</p></div>`)
		assert.Contains(t, string(p), `<input id="input" type="text" name="message" placeholder="Type your message here" required autofocus>`)
	})
}

func newRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	return req
}
