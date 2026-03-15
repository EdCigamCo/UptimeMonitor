package repository

import (
	"database/sql"
	"fmt"
	"time"
	"uptime_monitor/model"
)

// parseTimeFromDB parses time string from database to time.Time
// Tries RFC3339 format first, then falls back to "2006-01-02 15:04:05" format
func parseTimeFromDB(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02 15:04:05", timeStr)
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time '%s': %w", timeStr, err)
	}
	return parsedTime, nil
}

// CreateSite creates a new site in the database
func CreateSite(db *sql.DB, url string) (*model.Site, error) {
	query := `INSERT INTO sites (url, created_at) VALUES (?, ?)`

	createdAt := time.Now()
	// Use RFC3339 format to match what SQLite returns
	createdAtStr := createdAt.Format(time.RFC3339)
	result, err := db.Exec(query, url, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	site := &model.Site{
		ID:        id,
		URL:       url,
		CreatedAt: createdAt,
	}

	return site, nil
}

// GetAllSites retrieves all sites from the database
func GetAllSites(db *sql.DB) ([]model.Site, error) {
	query := `SELECT id, url, created_at FROM sites ORDER BY id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sites: %w", err)
	}
	defer rows.Close()

	var sites []model.Site
	for rows.Next() {
		var site model.Site
		var createdAtStr string

		if err := rows.Scan(&site.ID, &site.URL, &createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to scan site: %w", err)
		}

		// Parse created_at string to time.Time
		parsedTime, err := parseTimeFromDB(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		site.CreatedAt = parsedTime

		sites = append(sites, site)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sites: %w", err)
	}

	return sites, nil
}

// GetSiteByID retrieves a site by its ID
func GetSiteByID(db *sql.DB, id int64) (*model.Site, error) {
	query := `SELECT id, url, created_at FROM sites WHERE id = ?`

	var site model.Site
	var createdAtStr string

	err := db.QueryRow(query, id).Scan(&site.ID, &site.URL, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("site with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	// Parse created_at string to time.Time
	parsedTime, err := parseTimeFromDB(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	site.CreatedAt = parsedTime

	return &site, nil
}

// DeleteSite deletes a site from the database
func DeleteSite(db *sql.DB, id int64) error {
	query := `DELETE FROM sites WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("site with id %d not found", id)
	}

	return nil
}
