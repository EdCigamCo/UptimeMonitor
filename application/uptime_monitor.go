package application

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"uptime_monitor/infrastructure/repository"
	"uptime_monitor/model"
)

// ErrSiteNotFound is returned when a site is not found
var ErrSiteNotFound = errors.New("site not found")

// UptimeMonitor represents the application layer for uptime monitoring
type UptimeMonitor struct {
	db *sql.DB
}

// NewUptimeMonitor creates a new UptimeMonitor instance
func NewUptimeMonitor(db *sql.DB) *UptimeMonitor {
	return &UptimeMonitor{db: db}
}

// CreateSite creates a new site with URL validation
func (app *UptimeMonitor) CreateSite(siteURL string) (*model.Site, error) {
	if err := validateURL(siteURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	site, err := repository.CreateSite(app.db, siteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	return site, nil
}

// GetAllSites retrieves all sites
func (app *UptimeMonitor) GetAllSites() ([]model.Site, error) {
	sites, err := repository.GetAllSites(app.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get sites: %w", err)
	}

	return sites, nil
}

// DeleteSite deletes a site by ID
func (app *UptimeMonitor) DeleteSite(id int64) error {
	err := repository.DeleteSite(app.db, id)
	if err != nil {
		if errors.Is(err, repository.ErrSiteNotFound) {
			return fmt.Errorf("%w: site with id %d", ErrSiteNotFound, id)
		}
		return fmt.Errorf("failed to delete site: %w", err)
	}

	return nil
}

// ✅ Добавляем функцию валидации URL
// validateURL validates that the URL is properly formatted
func validateURL(s string) error {
	if s == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include host")
	}

	return nil
}
