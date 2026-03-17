package dto

// CreateSiteRequest represents request body for creating a new site
type CreateSiteRequest struct {
	URL string `json:"url"`
}

// SiteResponse represents a site in API response
type SiteResponse struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

// SiteListResponse represents list of sites in API response
type SiteListResponse struct {
	Sites []SiteResponse `json:"sites"`
}

// ErrorResponse represents error in API response
type ErrorResponse struct {
	Error string `json:"error"`
}
