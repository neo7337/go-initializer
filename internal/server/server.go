package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// validate is the shared request validator used by handlers.
var validate = validator.New()

// test seams for Start; defaults preserve production behavior.
var listenAndServe = func(s *http.Server) error { return s.ListenAndServe() }
var logFatalf = log.Fatalf

func Start() {
	routes := NewRouter(validate)

	server := &http.Server{
		Addr:           ":8182",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[INFO] Starting server on %s", server.Addr)
	err := listenAndServe(server)
	if err != nil {
		logFatalf("[ERROR] Server failed to start: %v", err)
	}
}
