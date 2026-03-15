package worker

import (
	"net/http"
	"time"
)

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
