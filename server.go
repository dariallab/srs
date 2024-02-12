package main

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `hello`)
}
