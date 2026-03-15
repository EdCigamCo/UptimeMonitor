package presentation

import (
	"fmt"
	"net/http"
	"time"
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
