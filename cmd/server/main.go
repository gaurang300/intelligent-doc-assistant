package main

import (
	"fmt"
	"log"
	"net/http"

	"intelligent-doc-assistant/api"
	"intelligent-doc-assistant/config"
)

func main() {
	cfg := config.GetConfig()

	// Initialize all services
	server := api.NewServer()

	// Start the HTTP server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
