package model

import "time"

// Site represents a website being monitored
type Site struct {
	ID        int64     `json:"id" db:"id"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Check represents a single availability check result
type Check struct {
	ID           int64     `json:"id" db:"id"`
	SiteID       int64     `json:"site_id" db:"site_id"`
	Status       string    `json:"status" db:"status"`
	ResponseTime int64     `json:"response_time" db:"response_time"`
	CheckedAt    time.Time `json:"checked_at" db:"checked_at"`
}
