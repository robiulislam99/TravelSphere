// services/attraction_service.go
// AttractionService fetches and transforms tourist attraction data
// from the OpenTripMap API. Falls back gracefully if the API is unavailable.
package services

import (
	"log"
	"strings"

	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/utils"
)

// AttractionService provides attraction lookup by geographic coordinates.
type AttractionService struct {
	client *utils.OpenTripMapClient
}

// NewAttractionService creates an AttractionService.
// Returns a service with a nil client if the API key is missing —
// all methods degrade gracefully returning empty slices instead of crashing.
func NewAttractionService() *AttractionService {
	client, err := utils.NewOpenTripMapClient()
	if err != nil {
		log.Printf("[AttractionService] WARNING: %v — attraction features disabled", err)
		return &AttractionService{client: nil}
	}
	return &AttractionService{client: client}
}

// GetByCoords fetches up to `limit` attractions near the given coordinates.
// radius is in meters. Returns an empty slice (not an error) on API failure
// so pages still render without attractions when the API is down.
func (s *AttractionService) GetByCoords(lat, lon float64, radius, limit int) ([]models.Attraction, error) {
	if s.client == nil {
		return []models.Attraction{}, nil
	}
	if radius <= 0 {
		radius = 10000 // default 10 km
	}
	if limit <= 0 {
		limit = 12
	}

	raw, err := s.client.GetAttractionsByRadius(lat, lon, radius, limit, "")
	if err != nil {
		log.Printf("[AttractionService] GetByCoords error: %v", err)
		return []models.Attraction{}, nil // graceful degradation
	}

	results := make([]models.Attraction, 0, len(raw.Features))
	for _, f := range raw.Features {
		name := strings.TrimSpace(f.Properties.Name)
		if name == "" {
			continue // skip unnamed attractions
		}
		lon, lat := 0.0, 0.0
		if len(f.Geometry.Coordinates) == 2 {
			lon = f.Geometry.Coordinates[0]
			lat = f.Geometry.Coordinates[1]
		}
		results = append(results, models.Attraction{
			XID:         f.Properties.XID,
			Name:        name,
			Kinds:       f.Properties.Kinds,
			PrimaryKind: utils.PrimaryKind(f.Properties.Kinds),
			Distance:    f.Properties.Dist,
			Latitude:    lat,
			Longitude:   lon,
		})
	}
	return results, nil
}

// GetForHomePage returns a small set of attractions across popular coordinates
// for display on the home page. Uses a few hardcoded world landmarks.
func (s *AttractionService) GetForHomePage() []models.Attraction {
	if s.client == nil {
		return []models.Attraction{}
	}
	// Paris — Eiffel Tower area
	attractions, _ := s.GetByCoords(48.8584, 2.2945, 5000, 4)
	return attractions
}