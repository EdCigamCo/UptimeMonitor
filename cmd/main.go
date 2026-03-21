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
	"uptime_monitor/application"
	"uptime_monitor/infrastructure/config"
	"uptime_monitor/infrastructure/worker"

	"golang.org/x/sync/errgroup"
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

	app := application.NewUptimeMonitor(db)
	handlers := presentation.NewHandlers(app)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)

	// Pass gCtx to worker so it stops when any goroutine fails
	checkInterval := 30 * time.Second
	w := worker.NewWorker(db, gCtx, checkInterval)

	port := cfg.Port

	// Register routes
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/info", handlers.InfoHandler)
	http.HandleFunc("/check", handlers.CheckHandler)
	// GET /api/sites - list all sites
	http.HandleFunc("/api/sites", handlers.ListSitesHandler)
	// POST /api/site - create site
	http.HandleFunc("/api/site", handlers.CreateSiteHandler)
	// DELETE /api/site/:id - delete site by ID
	http.HandleFunc("/api/site/", handlers.DeleteSiteHandler)
	// GET /api/sites/:id/history - get check history for a site
	http.HandleFunc("/api/sites/", handlers.GetSiteHistoryHandler)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server started on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET http://localhost:%s/health", port)
	log.Printf("  GET http://localhost:%s/info", port)
	log.Printf("  GET http://localhost:%s/check?url=<website_url>", port)
	log.Printf("  GET http://localhost:%s/api/sites", port)
	log.Printf("  POST http://localhost:%s/api/site", port)
	log.Printf("  DELETE http://localhost:%s/api/site/<id>", port)
	log.Printf("  GET http://localhost:%s/api/sites/<id>/history?limit=<n>", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{Addr: addr}

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		return w.Run()
	})

	g.Go(func() error {
		select {
		case <-sigChan:
			log.Println("Received shutdown signal")

			// Shutdown HTTP server first (before canceling context)
			log.Println("Initiating graceful shutdown...")
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			log.Println("Shutting down HTTP server...")
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
			} else {
				log.Println("HTTP server stopped gracefully")
			}

			// Cancel context to stop worker
			cancel()
			return nil
		case <-gCtx.Done():
			return gCtx.Err()
		}
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		log.Printf("Error in goroutine: %v", err)
	}

	// Worker should already be stopped via context cancellation
	log.Println("Worker stopped (context cancelled)")
	log.Println("Database connection closed")
	log.Println("Shutdown complete")
}
