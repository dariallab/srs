package main

import (
	"log"
	"net/http"
)

func main() {
	server := &Server{}

	if err := http.ListenAndServe(":8181", server); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
