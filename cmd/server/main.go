package main

import (
	"log"
	"net/http"

	"github.com/SkaMarius/go-agentic-base/internal/config"
	"github.com/SkaMarius/go-agentic-base/internal/server"
)

func main() {
	cfg := config.Load()
	router := server.NewRouter()

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
