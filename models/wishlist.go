// models/wishlist.go
// WishlistEntry is the core entity for the travel wishlist feature.
// It is stored in-memory (no database) and accessed via WishlistService.
package models

import (
	"errors"
	"strings"
	"time"
)

// WishlistStatus represents allowed travel status values.
type WishlistStatus string

const (
	StatusPlanned WishlistStatus = "Planned"
	StatusVisited WishlistStatus = "Visited"
)

// ValidStatuses is the set of allowed status values for validation.
var ValidStatuses = map[WishlistStatus]bool{
	StatusPlanned: true,
	StatusVisited: true,
}

// WishlistEntry represents one saved destination in a user's wishlist.
type WishlistEntry struct {
	ID          string         `json:"id"`
	CountryName string         `json:"country_name"`
	Slug        string         `json:"slug"`       // for linking to /countries/:slug
	Note        string         `json:"note"`       // optional user note
	Status      WishlistStatus `json:"status"`     // "Planned" or "Visited"
	CreatedAt   time.Time      `json:"created_at"` // set automatically on creation
}

// StatusClass returns a CSS class name for the status badge.
// Maps to .badge-planned and .badge-visited in main.css.
func (w *WishlistEntry) StatusClass() string {
	return strings.ToLower(string(w.Status))
}

// FormattedDate returns the creation date in a human-readable format.
func (w *WishlistEntry) FormattedDate() string {
	return w.CreatedAt.Format("2 Jan 2006")
}

// --- Request / Response DTOs ---

// CreateWishlistRequest is the JSON payload for POST /api/wishlist.
type CreateWishlistRequest struct {
	CountryName string `json:"country_name"`
	Note        string `json:"note"`
	Status      string `json:"status"` // must be "Planned" or "Visited"
}

// UpdateWishlistRequest is the JSON payload for PUT /api/wishlist/:id.
type UpdateWishlistRequest struct {
	Note   string `json:"note"`
	Status string `json:"status"`
}

// Validate checks CreateWishlistRequest fields and returns an error if invalid.
func (r *CreateWishlistRequest) Validate() error {
	if strings.TrimSpace(r.CountryName) == "" {
		return errors.New("country_name is required")
	}
	// Default status to Planned if not provided
	if r.Status == "" {
		r.Status = string(StatusPlanned)
	}
	if !ValidStatuses[WishlistStatus(r.Status)] {
		return errors.New("status must be 'Planned' or 'Visited'")
	}
	return nil
}

// Validate checks UpdateWishlistRequest fields and returns an error if invalid.
func (r *UpdateWishlistRequest) Validate() error {
	if r.Status == "" {
		return errors.New("status is required")
	}
	if !ValidStatuses[WishlistStatus(r.Status)] {
		return errors.New("status must be 'Planned' or 'Visited'")
	}
	return nil
}