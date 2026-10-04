package main

import (
	"log"
	"net/http"
	"time"

	"tasks-api/internal/handlers"
	httpmiddleware "tasks-api/internal/http"
	"tasks-api/internal/storage/memory"
)

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpmiddleware.Logging(handlers.New(memory.New())),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("server listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
