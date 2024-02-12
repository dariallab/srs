package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTutorServer(t *testing.T) {
	server := &Server{}

	t.Run("open page - see hello", func(t *testing.T) {
		req := newRequest(t, http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()
		want := `hello`

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, want, resp.Body.String())
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
