package srs

import (
	"net/http"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/auth"
	"github.com/dariallab/srs/pkg/templates/static"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

const errInternal = "There's an error, try again later"

type Server struct {
	http.Handler
	ai     ai.AI
	auth   auth.Auth
	store  sessions.Store
	logger zerolog.Logger
}

func NewServer(ai ai.AI, auth auth.Auth, store sessions.Store, logger zerolog.Logger) *Server {
	s := &Server{
		ai:     ai,
		auth:   auth,
		store:  store,
		logger: logger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.LoggedIn(s.showChatHandler))
	mux.HandleFunc("GET /message", s.LoggedIn(s.sendMessageHandler))
	mux.HandleFunc("GET /auth/login", s.authLoginHandler)
	mux.HandleFunc("GET /auth/logout", s.authLogoutHandler)
	mux.HandleFunc("POST /auth/callback", s.authCallbackHandler)

	staticServer := http.FileServer(http.FS(static.FS))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticServer))

	s.Handler = mux

	return s
}
