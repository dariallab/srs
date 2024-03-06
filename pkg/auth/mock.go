package auth

import (
	"context"
	"fmt"
	"net/http"
)

type Mock struct {
	AuthFn func(ctx context.Context, r *http.Request) (string, error)
}

func (m *Mock) Auth(ctx context.Context, r *http.Request) (string, error) {
	if m.AuthFn == nil {
		return "", fmt.Errorf("AuthFn is not implemented")
	}
	return m.AuthFn(ctx, r)
}

func (m *Mock) GetClientID() string {
	return "mock-client-id"
}

func (m *Mock) GetCallbackURL() string {
	return "mock-callback-url"
}
