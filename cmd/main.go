package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	http.HandleFunc("/check", presentation.CheckHandler)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server started on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET http://localhost:%s/health", port)
	log.Printf("  GET http://localhost:%s/info", port)
	log.Printf("  GET http://localhost:%s/check?url=<website_url>", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{Addr: addr}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
