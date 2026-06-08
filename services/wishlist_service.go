// services/wishlist_service.go
// WishlistService manages wishlist entries in an in-memory store protected
// by a mutex. No database is used — this satisfies the assessment requirement.
// All CRUD operations are exposed as methods; controllers never touch storage directly.
package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/utils"
)

// WishlistService is safe for concurrent use.
type WishlistService struct {
	mu      sync.RWMutex
	entries map[string]*models.WishlistEntry // keyed by ID
}

// NewWishlistService creates an empty in-memory wishlist store.
func NewWishlistService() *WishlistService {
	return &WishlistService{
		entries: make(map[string]*models.WishlistEntry),
	}
}

// GetAll returns all wishlist entries as a slice, sorted by creation time (newest first).
func (s *WishlistService) GetAll() []models.WishlistEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]models.WishlistEntry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, *e)
	}
	// Sort newest first
	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].CreatedAt.After(list[i].CreatedAt) {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	return list
}

// Create validates and adds a new wishlist entry.
// Returns the created entry or a validation error.
func (s *WishlistService) Create(req *models.CreateWishlistRequest) (*models.WishlistEntry, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	entry := &models.WishlistEntry{
		ID:          uuid.New().String(),
		CountryName: utils.SanitizeString(req.CountryName),
		Slug:        utils.NameToSlug(req.CountryName),
		Note:        utils.SanitizeString(req.Note),
		Status:      models.WishlistStatus(req.Status),
		CreatedAt:   time.Now(),
	}

	s.mu.Lock()
	s.entries[entry.ID] = entry
	s.mu.Unlock()

	return entry, nil
}

// GetByID returns a single wishlist entry by ID.
// Returns nil, nil when not found.
func (s *WishlistService) GetByID(id string) (*models.WishlistEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.entries[id]
	if !ok {
		return nil, nil
	}
	copy := *entry
	return &copy, nil
}

// Update applies note and status changes to an existing entry.
// Returns the updated entry or an error if not found / invalid.
func (s *WishlistService) Update(id string, req *models.UpdateWishlistRequest) (*models.WishlistEntry, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.entries[id]
	if !ok {
		return nil, fmt.Errorf("wishlist entry not found")
	}

	entry.Note = utils.SanitizeString(req.Note)
	entry.Status = models.WishlistStatus(req.Status)

	copy := *entry
	return &copy, nil
}

// Delete removes a wishlist entry by ID.
// Returns an error if the entry does not exist.
func (s *WishlistService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[id]; !ok {
		return fmt.Errorf("wishlist entry not found")
	}
	delete(s.entries, id)
	return nil
}

// IsDuplicate checks whether a country is already in the wishlist (case-insensitive).
// Useful for showing a warning before adding a duplicate.
func (s *WishlistService) IsDuplicate(countryName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	name := strings.ToLower(strings.TrimSpace(countryName))
	for _, e := range s.entries {
		if strings.ToLower(e.CountryName) == name {
			return true
		}
	}
	return false
}