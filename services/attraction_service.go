// services/attraction_service.go
package services

import (
	"log"
	"strings"

	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/utils"
)

type AttractionService struct {
	client *utils.OpenTripMapClient
}

func NewAttractionService() *AttractionService {
	client, err := utils.NewOpenTripMapClient()
	if err != nil {
		log.Printf("[AttractionService] WARNING: %v — attraction features disabled", err)
		return &AttractionService{client: nil}
	}
	return &AttractionService{client: client}
}

func (s *AttractionService) GetByCoords(lat, lon float64, radius, limit int) ([]models.Attraction, error) {
	if s.client == nil {
		return []models.Attraction{}, nil
	}
	if radius <= 0 {
		radius = 10000
	}
	if limit <= 0 {
		limit = 12
	}

	raw, err := s.client.GetAttractionsByRadius(lat, lon, radius, limit, "")
	if err != nil {
		log.Printf("[AttractionService] GetByCoords error: %v", err)
		return []models.Attraction{}, nil
	}

	results := make([]models.Attraction, 0, len(raw.Features))
	for _, f := range raw.Features {
		name := strings.TrimSpace(f.Properties.Name)
		if name == "" {
			continue
		}

		// Use var to avoid shadowing the outer lat/lon params
		var featLon, featLat float64
		if len(f.Geometry.Coordinates) == 2 {
			featLon = f.Geometry.Coordinates[0] // GeoJSON: [lon, lat]
			featLat = f.Geometry.Coordinates[1]
		}

		results = append(results, models.Attraction{
			XID:         f.Properties.XID,
			Name:        name,
			Kinds:       f.Properties.Kinds,
			PrimaryKind: utils.PrimaryKind(f.Properties.Kinds),
			Distance:    f.Properties.Dist,
			Latitude:    featLat,
			Longitude:   featLon,
		})
	}
	return results, nil
}

// homePageSeeds are fallback locations tried in order until one returns results.
var homePageSeeds = []struct {
	label    string
	lat, lon float64
}{
	{"Paris", 48.8584, 2.2945},
	{"Tokyo", 35.6762, 139.6503},
	{"New York", 40.7128, -74.0060},
}

func (s *AttractionService) GetForHomePage() []models.Attraction {
	if s.client == nil {
		return []models.Attraction{}
	}

	for _, seed := range homePageSeeds {
		attractions, err := s.GetByCoords(seed.lat, seed.lon, 5000, 4)
		if err != nil {
			log.Printf("[AttractionService] GetForHomePage error for %s: %v", seed.label, err)
			continue
		}
		if len(attractions) > 0 {
			return attractions
		}
		log.Printf("[AttractionService] GetForHomePage: no results for %s, trying next seed", seed.label)
	}

	log.Printf("[AttractionService] GetForHomePage: all seeds exhausted, returning empty")
	return []models.Attraction{}
}
