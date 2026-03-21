package dto

// CreateSiteRequest represents request body for creating a new site
type CreateSiteRequest struct {
	URL string `json:"url"`
}

// SiteResponse represents a site in API response
type SiteResponse struct {
	ID           int64  `json:"id"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
	Status       string `json:"status,omitempty"`
	ResponseTime int64  `json:"response_time,omitempty"`
	LastChecked  string `json:"last_checked,omitempty"`
}

// SiteListResponse represents list of sites in API response
type SiteListResponse struct {
	Sites []SiteResponse `json:"sites"`
}

// ErrorResponse represents error in API response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CheckResponse represents a single check in API response
type CheckResponse struct {
	ID           int64  `json:"id"`
	Status       string `json:"status"`        // "up" or "down"
	ResponseTime int64  `json:"response_time"` // in milliseconds
	CheckedAt    string `json:"checked_at"`    // RFC3339 format
}

// SiteHistoryResponse represents check history for a site
type SiteHistoryResponse struct {
	SiteID int64           `json:"site_id"`
	URL    string          `json:"url"`
	Checks []CheckResponse `json:"checks"`
}
