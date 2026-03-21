package worker

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"uptime_monitor/infrastructure/repository"
)

// Worker represents a background worker for checking site availability
type Worker struct {
	db            *sql.DB
	checkInterval time.Duration
	ctx           context.Context
}

// NewWorker creates a new worker instance
func NewWorker(db *sql.DB, ctx context.Context, checkInterval time.Duration) *Worker {
	return &Worker{
		db:            db,
		checkInterval: checkInterval,
		ctx:           ctx,
	}
}

// Run is the main worker loop
// Returns error if context is cancelled or if there's a critical error
func (w *Worker) Run() error {
	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()

	log.Printf("Worker started with check interval: %v", w.checkInterval)

	// Perform initial check immediately
	w.checkAllSites()

	for {
		select {
		case <-w.ctx.Done():
			log.Println("Worker stopped")
			return w.ctx.Err()
		case <-ticker.C:
			w.checkAllSites()
		}
	}
}

// checkAllSites checks all sites and saves results to database
func (w *Worker) checkAllSites() {
	// Get all sites from database
	sites, err := repository.GetAllSites(w.db)
	if err != nil {
		log.Printf("Error getting sites for check: %v", err)
		return
	}

	if len(sites) == 0 {
		log.Println("No sites to check")
		return
	}

	log.Printf("Checking %d sites...", len(sites))

	// Check each site
	for _, site := range sites {
		// Check site availability
		status, responseTime, err := CheckSiteAvailability(site.URL)
		if err != nil {
			log.Printf("Error checking site %s (ID: %d): %v", site.URL, site.ID, err)
			// Save check result even if there was an error (status will be "down")
			status = "down"
			responseTime = 0
		}

		// Save check result to database
		if err := repository.SaveCheck(w.db, site.ID, status, responseTime); err != nil {
			log.Printf("Error saving check result for site %s (ID: %d): %v", site.URL, site.ID, err)
			continue
		}

		log.Printf("Checked site %s (ID: %d): %s, response time: %d ms", site.URL, site.ID, status, responseTime)
	}

	log.Println("Finished checking all sites")
}

// CheckSiteAvailability checks if a website is available
// Returns status ("up" or "down"), response time in milliseconds, and error if any
func CheckSiteAvailability(url string) (status string, responseTime int64, err error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Record start time
	start := time.Now()

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		return "down", 0, err
	}
	defer resp.Body.Close()

	// Calculate response time
	responseTime = time.Since(start).Milliseconds()

	// Determine status based on HTTP status code
	// Status codes 200-299 are considered "up"
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status = "up"
	} else {
		status = "down"
	}

	return status, responseTime, nil
}
