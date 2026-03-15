package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"uptime_monitor/presentation"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get port from environment variable or use default value
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Register routes
	http.HandleFunc("/health", presentation.HealthHandler)
	http.HandleFunc("/info", presentation.InfoHandler)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server started on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET http://localhost:%s/health", port)
	log.Printf("  GET http://localhost:%s/info", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
