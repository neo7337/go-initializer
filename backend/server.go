package main

import (
	"net/http"
	"time"

	"log"

	"github.com/go-playground/validator/v10"
)

func Start() {
	validate := validator.New()

	routes := NewRouter(validate)

	server := &http.Server{
		Addr:           ":8182",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[INFO] Starting server on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("[ERROR] Server failed to start: %v", err)
	}
}
