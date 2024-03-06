package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dariallab/srs/pkg/ai"
	"github.com/dariallab/srs/pkg/srs"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        int    `envconfig:"PORT" required:"true"`
	OpenAIToken string `envconfig:"OPENAI_TOKEN" required:"true"`
}

func main() {
	l := zerolog.New(os.Stdout)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		l.Fatal().Err(err).Msg("can't load confg")
	}

	// ai := ai.New(cfg.OpenAIToken)
	ai := ai.NewMock()

	server := srs.NewServer(ai, l)
	l.Info().Int("port", cfg.Port).Msg("starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), server); err != nil {
		l.Fatal().Err(err).Msg("can't start server")
	}
}
