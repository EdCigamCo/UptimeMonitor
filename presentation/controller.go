package presentation

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"uptime_monitor/infrastructure/worker"
)

// HealthHandler handles GET /health request
// Returns simple text "OK" for server health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// InfoHandler handles GET /info request
// Returns simple text with server information
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	info := fmt.Sprintf("UptimeMonitor Server\nLast request at: %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprint(w, info)
}

// CheckHandler handles GET /check request
// Checks availability of a website specified in query parameter "url"
// Returns simple text with status and response time
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get URL from query parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	// Check site availability
	status, responseTime, err := worker.CheckSiteAvailability(url)
	if err != nil {
		log.Printf("Error checking site %s: %v", url, err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "URL: %s\nStatus: down\nError: %v", url, err)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "URL: %s\nStatus: %s\nResponse time: %d ms", url, status, responseTime)
}
