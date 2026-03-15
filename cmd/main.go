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
	"uptime_monitor/infrastructure/config"

	"uptime_monitor/presentation"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := config.InitDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDatabase(db)

	log.Println("Database initialized successfully")

	handlers := presentation.NewHandlers()
	port := cfg.Port

	// Register routes
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/info", handlers.InfoHandler)
	http.HandleFunc("/check", handlers.CheckHandler)

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
	log.Println("Database connection closed")
}
