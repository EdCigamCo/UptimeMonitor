package presentation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"uptime_monitor/application"
	"uptime_monitor/application/dto"
	"uptime_monitor/infrastructure/worker"
	"uptime_monitor/presentation/mapper"
)

// Handlers holds application layer dependencies
type Handlers struct {
	app *application.UptimeMonitor
}

func NewHandlers(app *application.UptimeMonitor) *Handlers {
	return &Handlers{app: app}
}

// HealthHandler handles GET /health request
// Returns simple text "OK" for server health check
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
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
func (h *Handlers) InfoHandler(w http.ResponseWriter, r *http.Request) {
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
func (h *Handlers) CheckHandler(w http.ResponseWriter, r *http.Request) {
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

// CreateSiteHandler handles POST /api/site
func (h *Handlers) CreateSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req dto.CreateSiteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	site, err := h.app.CreateSite(req.URL)
	if err != nil {
		log.Printf("Error creating site: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := mapper.ToSiteResponse(site)

	respondWithJSON(w, http.StatusCreated, response)
}

// ListSitesHandler handles GET /api/sites
func (h *Handlers) ListSitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sites, err := h.app.GetAllSites()
	if err != nil {
		log.Printf("Error getting sites: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve sites")
		return
	}

	response := mapper.ToSiteListResponse(sites)

	respondWithJSON(w, http.StatusOK, response)
}

// DeleteSiteHandler handles DELETE /api/site/:id
func (h *Handlers) DeleteSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/site/")
	if path == "" {
		respondWithError(w, http.StatusBadRequest, "Site ID is required")
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	if err := h.app.DeleteSite(id); err != nil {
		if errors.Is(err, application.ErrSiteNotFound) {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Site with id %d not found", id))
			return
		}
		log.Printf("Error deleting site: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete site")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondWithJSON sends JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondWithError sends error response as JSON
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := dto.ErrorResponse{
		Error: message,
	}
	respondWithJSON(w, statusCode, response)
}
