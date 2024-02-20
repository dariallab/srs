package srs

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestTutorServer(t *testing.T) {

	t.Run("open page - see form", func(t *testing.T) {
		server := NewServer(&ai.Mock{}, zerolog.New(os.Stdout))
		req := newRequest(t, http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), `<form ws-send`, "websocket tags are not found")
	})

	t.Run("send message - return response via websockets", func(t *testing.T) {
		server := NewServer(&ai.Mock{
			CorrectionFn: func(s string) (string, error) {
				return "hallo", nil
			},
		}, zerolog.New(os.Stdout))
		srv, ws := wsWriteMessage(t, server, "/message", "hello")
		defer ws.Close()
		defer srv.Close()

		got := wsReadMessage(t, ws)
		wantOriginalMsg := `<textarea id="chat_input" name="message" oninput="this.style.height = ''; this.style.height = this.scrollHeight +'px'" rows="1" required autofocus placeholder="Type your message here"class="flex-1 mr-2 resize-none bg-slate-300 focus:border-none focus:outline-none"></textarea><div id="chat_messages" hx-swap-oob="afterbegin" class="h-full flex flex-col-reverse overflow-auto flex-1 p-6 mb-auto"><p>hello</p></div>`
		assert.Equal(t, wantOriginalMsg, got)

		got = wsReadMessage(t, ws)
		wantDiffMsg := `<div id="chat_messages" hx-swap-oob="afterbegin" class="h-full flex flex-col-reverse overflow-auto flex-1 p-6 mb-auto"><p>h<span class="bg-red-200">e</span>llo</p><p>h<span class="bg-green-200">a</span>llo</p></div>`
		assert.Equal(t, wantDiffMsg, got)
	})

	t.Run("serve static", func(t *testing.T) {
		server := NewServer(&ai.Mock{}, zerolog.New(os.Stdout))
		req := newRequest(t, http.MethodGet, "/src/input.css", nil)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func wsWriteMessage(t *testing.T, server http.Handler, endpoint, message string) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	srv := httptest.NewServer(server)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + endpoint

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("can't open a ws connection on %s %v", wsURL, err)
	}

	msg := fmt.Sprintf(`{"message":"%s"}`, message)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("can't send message over ws connection %v", err)
	}

	return srv, ws
}

func wsReadMessage(t *testing.T, ws *websocket.Conn) string {
	t.Helper()

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("can't read the websocket message: %s ", err)
	}

	got := strings.ReplaceAll(string(p), "\n", "")
	got = strings.ReplaceAll(got, "\t", "")
	return got
}

func newRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	return req
}
