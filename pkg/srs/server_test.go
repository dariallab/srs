package srs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTutorServer(t *testing.T) {
	server := Server{}

	t.Run("open page - see form", func(t *testing.T) {
		req := newRequest(t, http.MethodGet, "/chat", nil)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), `<form hx-post="/chat"`, "form not found")
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
