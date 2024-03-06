package srs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/auth"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Run("show login page", func(t *testing.T) {
		server := NewServer(&ai.Mock{}, &auth.Mock{}, sessions.NewCookieStore(), zerolog.New(os.Stdout))
		req := newRequest(t, server, http.MethodGet, "/auth/login", nil, false)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), `<div class="g_id_signin"`, "don't contain sign-in button")
	})

	t.Run("not logged in - redirect to login page", func(t *testing.T) {
		server := NewServer(&ai.Mock{}, &auth.Mock{}, sessions.NewCookieStore(), zerolog.New(os.Stdout))
		req := newRequest(t, server, http.MethodGet, "/", nil, false)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
		assert.Contains(t, resp.Body.String(), `/auth/login`, "didn't redirect")
	})
	t.Run("logged in - see chat page", func(t *testing.T) {
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
		chatReq := newRequest(t, server, http.MethodGet, "/", nil, true)
		chatResp := httptest.NewRecorder()

		server.ServeHTTP(chatResp, chatReq)

		assert.Equal(t, http.StatusOK, chatResp.Code)
		assert.Contains(t, chatResp.Body.String(), `<div hx-ext="ws"`, "don't contain websocket")
	})

	t.Run("logged out - redirect to login page", func(t *testing.T) {
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
		logoutReq := newRequest(t, server, http.MethodGet, "/auth/logout", nil, true)
		logoutResp := httptest.NewRecorder()

		server.ServeHTTP(logoutResp, logoutReq)

		assert.Equal(t, http.StatusTemporaryRedirect, logoutResp.Code)

		chatReq := newRequest(t, server, http.MethodGet, "/", nil, false)
		chatResp := httptest.NewRecorder()

		server.ServeHTTP(chatResp, chatReq)

		assert.Equal(t, http.StatusTemporaryRedirect, chatResp.Code)
		assert.Contains(t, chatResp.Body.String(), `/auth/login`, "didn't redirect")
	})
}
