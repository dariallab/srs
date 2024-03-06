package srs

import (
	"net/http"

	"github.com/dariallab/srs/pkg/templates"
)

type templateDataLogin struct {
	ClientID    string
	CallbackURL string
}

func (s *Server) authLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t := templates.TemplateLogin
	data := &templateDataLogin{
		ClientID:    s.auth.GetClientID(),
		CallbackURL: s.auth.GetCallbackURL(),
	}

	if err := t.Execute(w, data); err != nil {
		s.logger.Error().Err(err).Msg("can't execute template")
		http.Error(w, errInternal, http.StatusInternalServerError)
		return
	}
}

func (s *Server) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authID, err := s.auth.Auth(ctx, r)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get auth")
		http.Error(w, errInternal, http.StatusInternalServerError)
		return
	}

	if err := s.saveUserToSession(r, w, authID); err != nil {
		s.logger.Error().Err(err).Msg("can't set session")
		http.Error(w, errInternal, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) authLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.cleanUserSession(w, r); err != nil {
		s.logger.Error().Err(err).Msg("can't clean session")
		http.Error(w, errInternal, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
