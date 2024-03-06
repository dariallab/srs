package srs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/auth"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestChatServer(t *testing.T) {

	t.Run("logged in, open chat - see form", func(t *testing.T) {
		server := NewServer(
			&ai.Mock{},
			&auth.Mock{
				AuthFn: func(ctx context.Context, r *http.Request) (string, error) {
					return "test-id", nil
				},
			},
			sessions.NewCookieStore([]byte("test-key")),
			zerolog.New(os.Stdout),
		)
		req := newRequest(t, server, http.MethodGet, "/", nil, true)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), `<form ws-send`, "websocket tags are not found")
	})

	t.Run("logged in, send message - return response via websockets", func(t *testing.T) {
		server := NewServer(
			&ai.Mock{
				CorrectionFn: func(s string) (string, error) {
					return "hallo", nil
				},
				ResponseFn: func(s string) (string, error) {
					return "response", nil
				},
			},
			&auth.Mock{
				AuthFn: func(ctx context.Context, r *http.Request) (string, error) {
					return "test-id", nil
				},
			},
			sessions.NewCookieStore([]byte("test-key")),
			zerolog.New(os.Stdout),
		)
		srv, ws := wsWriteMessage(t, server, "/message", "hello", true)
		defer ws.Close()
		defer srv.Close()

		got := wsReadMessage(t, ws)
		assert.Contains(t, got, `hello`, "original message is not found")

		got = wsReadMessage(t, ws)
		assert.Contains(t, got, `<span class="bg-red-100">e</span>`, "diff message is not found")

		got = wsReadMessage(t, ws)
		assert.Contains(t, got, `response`, "response message is not found")
	})

}

func wsWriteMessage(t *testing.T, server http.Handler, endpoint, message string, loggedIn bool) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	srv := httptest.NewServer(server)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + endpoint

	dialer := websocket.Dialer{}
	header := http.Header{}
	if loggedIn {
		cookies := getAuthCookies(t, server)
		for _, cookie := range cookies {
			header.Add("Cookie", cookie.String())
		}
	}
	ws, _, err := dialer.Dial(wsURL, header)
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
