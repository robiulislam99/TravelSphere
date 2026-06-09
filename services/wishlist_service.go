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

type WishlistService struct {
    mu      sync.RWMutex
    entries map[string]map[string]*models.WishlistEntry // username → id → entry
}

func NewWishlistService() *WishlistService {
    return &WishlistService{
        entries: make(map[string]map[string]*models.WishlistEntry),
    }
}

// userEntries returns (creating if needed) the entry map for a user.
// Must be called with lock held.
func (s *WishlistService) userEntries(username string) map[string]*models.WishlistEntry {
    if _, ok := s.entries[username]; !ok {
        s.entries[username] = make(map[string]*models.WishlistEntry)
    }
    return s.entries[username]
}

func (s *WishlistService) GetAll(username string) []models.WishlistEntry {
    s.mu.RLock()
    defer s.mu.RUnlock()

    list := make([]models.WishlistEntry, 0)
    for _, e := range s.userEntries(username) {
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

func (s *WishlistService) Create(username string, req *models.CreateWishlistRequest) (*models.WishlistEntry, error) {
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
    s.userEntries(username)[entry.ID] = entry
    s.mu.Unlock()

    return entry, nil
}

func (s *WishlistService) GetByID(username, id string) (*models.WishlistEntry, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    entry, ok := s.userEntries(username)[id]
    if !ok {
        return nil, nil
    }
    copy := *entry
    return &copy, nil
}

func (s *WishlistService) Update(username, id string, req *models.UpdateWishlistRequest) (*models.WishlistEntry, error) {
    if err := req.Validate(); err != nil {
        return nil, err
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    entry, ok := s.userEntries(username)[id]
    if !ok {
        return nil, fmt.Errorf("wishlist entry not found")
    }

    entry.Note   = utils.SanitizeString(req.Note)
    entry.Status = models.WishlistStatus(req.Status)
    copy := *entry
    return &copy, nil
}

func (s *WishlistService) Delete(username, id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.userEntries(username)[id]; !ok {
        return fmt.Errorf("wishlist entry not found")
    }
    delete(s.userEntries(username), id)
    return nil
}

func (s *WishlistService) IsDuplicate(username, countryName string) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()

    name := strings.ToLower(strings.TrimSpace(countryName))
    for _, e := range s.userEntries(username) {
        if strings.ToLower(e.CountryName) == name {
            return true
        }
    }
    return false
}