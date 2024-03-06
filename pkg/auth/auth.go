package auth

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"google.golang.org/api/idtoken"
)

type Auth interface {
	Auth(ctx context.Context, r *http.Request) (string, error)
	GetClientID() string
	GetCallbackURL() string
}

type Client struct {
	clientID    string
	callbackURL string
}

func New(clientID, callbackURL string) *Client {
	return &Client{
		clientID:    clientID,
		callbackURL: callbackURL,
	}
}

func (c *Client) Auth(ctx context.Context, r *http.Request) (string, error) {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read request body")
	}
	values, err := url.ParseQuery(string(all))
	if err != nil {
		return "", errors.Wrap(err, "can't parse query")
	}
	credential := values.Get("credential")

	csrfTokenBody := values.Get("g_csrf_token")
	if csrfTokenBody == "" {
		return "", errors.Wrap(err, "can't get g_csrf_token from data")
	}
	csrfCookie, err := r.Cookie("g_csrf_token")
	if err != nil {
		return "", errors.Wrap(err, "can't get g_csrf_token from cookie")
	}
	if csrfCookie.Value != csrfTokenBody {
		return "", errors.Wrap(err, "g_csrf_tokens don't match")
	}

	payload, err := idtoken.Validate(ctx, credential, c.clientID)
	if err != nil {
		return "", errors.Wrap(err, "invalid idtoken")
	}

	return payload.Subject, nil
}

func (c *Client) GetClientID() string {
	return c.clientID
}

func (c *Client) GetCallbackURL() string {
	return c.callbackURL
}
