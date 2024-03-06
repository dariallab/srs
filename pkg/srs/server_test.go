package srs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/auth"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestStaticServer(t *testing.T) {
	t.Run("serve static", func(t *testing.T) {
		server := NewServer(&ai.Mock{}, &auth.Mock{}, sessions.NewCookieStore(), zerolog.New(os.Stdout))
		req := newRequest(t, server, http.MethodGet, "/static/output.css", nil, false)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func newRequest(t *testing.T, server http.Handler, method, path string, body io.Reader, loggedIn bool) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	if loggedIn {
		cookies := getAuthCookies(t, server)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	return req
}

func getAuthCookies(t *testing.T, server http.Handler) []*http.Cookie {
	t.Helper()

	loginReq := newRequest(t, server, http.MethodPost, "/auth/callback", strings.NewReader(""), false)
	loginResp := httptest.NewRecorder()

	server.ServeHTTP(loginResp, loginReq)
	assert.Equal(t, http.StatusSeeOther, loginResp.Code)

	cookies := loginResp.Result().Cookies()
	assert.NotEmpty(t, cookies)
	return cookies
}
