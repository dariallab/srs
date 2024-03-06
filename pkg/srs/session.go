package srs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	sessionCookie    = "srs-sid"
	sessionUserIDKey = "user-id"
)

var contextKeyUser = contextKey("user")

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

type contextUserData struct {
	ID string
}

func getUserFromContext(ctx context.Context) (*contextUserData, error) {
	user, ok := ctx.Value(contextKeyUser).(*contextUserData)
	if !ok {
		return nil, fmt.Errorf("context contains invalid user")
	}
	return user, nil
}

func (s *Server) saveUserToContext(ctx context.Context, r *http.Request) (context.Context, error) {
	session, err := s.store.Get(r, sessionCookie)
	if err != nil {
		return ctx, err
	}
	if userID, ok := session.Values[sessionUserIDKey].(string); ok && userID != "" {
		ctx = context.WithValue(ctx, contextKeyUser, &contextUserData{ID: userID})
	}
	return ctx, err
}

func (s *Server) saveUserToSession(r *http.Request, w http.ResponseWriter, userID string) error {
	session, err := s.store.Get(r, sessionCookie)
	if err != nil {
		return err
	}
	session.Values[sessionUserIDKey] = userID
	return session.Save(r, w)
}

func (s *Server) cleanUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, sessionCookie)
	if err != nil {
		return errors.Wrap(err, "can't get session from store")
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return errors.Wrap(err, "can't save session")
	}

	return nil
}

func (s *Server) LoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := s.saveUserToContext(r.Context(), r)
		if err != nil {
			s.logger.Error().Err(err).Msg("can't set context from session")
			http.Error(w, errInternal, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
