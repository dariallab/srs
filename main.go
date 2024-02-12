package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int `envconfig:"PORT" required:"true" default:"8181"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("can't load confg: %v", err)
	}

	server := &Server{}

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), server); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
