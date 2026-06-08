// services/country_service.go
// CountryService orchestrates REST Countries API calls and transforms
// raw responses into application models. Controllers never touch the API
// client directly — all logic lives here.
package services

import (
	"sort"
	"strings"

	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/utils"
)

// CountryService provides country search and lookup operations.
type CountryService struct {
	client *utils.RestCountriesClient
}

// NewCountryService creates a CountryService with a REST Countries client.
func NewCountryService() *CountryService {
	return &CountryService{client: utils.NewRestCountriesClient()}
}

// GetAll returns all countries, optionally filtered by search query and/or region.
// Both filters are case-insensitive. Empty strings mean "no filter".
func (s *CountryService) GetAll(search, region string) ([]models.CountryListItem, error) {
	raw, err := s.client.GetAll()
	if err != nil {
		return nil, err
	}

	search = strings.ToLower(strings.TrimSpace(search))
	region = strings.TrimSpace(region)

	results := make([]models.CountryListItem, 0, len(raw))
	for _, r := range raw {
		// Region filter (exact match)
		if region != "" && !strings.EqualFold(r.Region, region) {
			continue
		}
		// Search filter (name contains query)
		if search != "" && !strings.Contains(strings.ToLower(r.Name.Common), search) {
			continue
		}
		results = append(results, s.toListItem(r))
	}

	// Sort alphabetically by name for consistent ordering
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results, nil
}

// GetBySlug finds a single country whose slug matches the provided value.
// Slug format: lowercase name with spaces replaced by hyphens e.g. "united-states".
// Returns nil, nil when not found (controller renders 404).
func (s *CountryService) GetBySlug(slug string) (*models.Country, error) {
	raw, err := s.client.GetAll()
	if err != nil {
		return nil, err
	}

	slug = strings.ToLower(strings.TrimSpace(slug))
	for _, r := range raw {
		if utils.NameToSlug(r.Name.Common) == slug {
			c := s.toCountry(r)
			return &c, nil
		}
	}
	return nil, nil // not found — caller handles 404
}

// GetFeatured returns a fixed set of popular countries for the home page.
// Fetches all countries and picks a curated subset by CCA2 code.
func (s *CountryService) GetFeatured() ([]models.CountryListItem, error) {
	featured := []string{"JP", "FR", "BR", "IN", "EG", "AU", "BD", "CA", "IT", "ZA"}
	raw, err := s.client.GetAll()
	if err != nil {
		return nil, err
	}

	index := make(map[string]models.CountryListItem, len(raw))
	for _, r := range raw {
		index[r.CCA2] = s.toListItem(r)
	}

	results := make([]models.CountryListItem, 0, len(featured))
	for _, code := range featured {
		if item, ok := index[code]; ok {
			results = append(results, item)
		}
	}
	return results, nil
}

// --- Private transformation helpers ---

// toCountry converts a RawCountry into a full models.Country.
func (s *CountryService) toCountry(r utils.RawCountry) models.Country {
	capital := ""
	if len(r.Capital) > 0 {
		capital = r.Capital[0]
	}

	lat, lon := 0.0, 0.0
	if len(r.Latlng) == 2 {
		lat, lon = r.Latlng[0], r.Latlng[1]
	}

	return models.Country{
		Name:                r.Name.Common,
		Slug:                utils.NameToSlug(r.Name.Common),
		CCA2:                r.CCA2,
		CCA3:                r.CCA3,
		Capital:             capital,
		Region:              r.Region,
		Subregion:           r.Subregion,
		Population:          r.Population,
		FormattedPopulation: utils.FormatPopulation(r.Population),
		FlagURL:             r.Flags.PNG,
		FlagEmoji:           r.Flag,
		CurrencyDisplay:     utils.FormatCurrencies(r.Currencies),
		LanguageDisplay:     utils.FormatLanguages(r.Languages),
		Latitude:            lat,
		Longitude:           lon,
	}
}

// toListItem converts a RawCountry into a lightweight CountryListItem.
func (s *CountryService) toListItem(r utils.RawCountry) models.CountryListItem {
	c := s.toCountry(r)
	return c.ToListItem()
}