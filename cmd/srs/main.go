package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/auth"
	"github.com/dariallab/srs/pkg/srs"
	"github.com/gorilla/sessions"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

func main() {
	l := zerolog.New(os.Stdout)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		l.Fatal().Err(err).Msg("can't load confg")
	}

	ai := ai.New(cfg.OpenAIToken)

	auth := auth.New(cfg.Auth.ClientID, cfg.Auth.CallbackURL)

	store := sessions.NewCookieStore([]byte(cfg.Auth.SessionKey))

	server := srs.NewServer(ai, auth, store, l)
	l.Info().Int("port", cfg.Port).Msg("starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), server); err != nil {
		l.Fatal().Err(err).Msg("can't start server")
	}
}
