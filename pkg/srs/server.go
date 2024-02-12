package srs

import (
	"net/http"

	"github.com/dariallab/srs/pkg/static"
	"github.com/rs/zerolog"
)

type Server struct {
	logger zerolog.Logger
}

func NewServer(logger zerolog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t := static.TemplateChat
	if err := t.Execute(w, nil); err != nil {
		s.logger.Error().Err(err).Msg("can't execute template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
