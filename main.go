package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	Port int `envconfig:"PORT" required:"true"`
}

func main() {
	l := zerolog.New(os.Stdout)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		l.Fatal().Err(err).Msg("can't load confg")
	}

	server := &Server{}
	l.Info().Int("port", cfg.Port).Msg("starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), server); err != nil {
		l.Fatal().Err(err).Msg("can't start server")
	}
}
