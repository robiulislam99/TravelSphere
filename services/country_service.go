// services/country_service.go
// CountryService orchestrates REST Countries v5 API calls and transforms
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

// NewCountryServiceWithClient creates a CountryService with a custom client.
// Used in tests to inject a mock REST Countries client.
func NewCountryServiceWithClient(client *utils.RestCountriesClient) *CountryService {
	return &CountryService{client: client}
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
		if region != "" && !strings.EqualFold(r.Region(), region) {
			continue
		}
		// Search filter (name contains query)
		if search != "" && !strings.Contains(strings.ToLower(r.Name()), search) {
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
		if utils.NameToSlug(r.Name()) == slug {
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
		index[r.CCA2()] = s.toListItem(r)
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

// toCountry converts a RawCountry (v5 map-based) into a full models.Country.
func (s *CountryService) toCountry(r utils.RawCountry) models.Country {
	// Build currency display string e.g. "BDT (Bangladeshi taka), USD"
	currencies := r.Currencies()
	currencyParts := make([]string, 0, len(currencies))
	for code, name := range currencies {
		if name != "" && name != code {
			currencyParts = append(currencyParts, code+" ("+name+")")
		} else {
			currencyParts = append(currencyParts, code)
		}
	}
	currencyDisplay := "N/A"
	if len(currencyParts) > 0 {
		currencyDisplay = strings.Join(currencyParts, ", ")
	}

	// Build language display string e.g. "Bengali, English"
	languages := r.Languages()
	languageDisplay := "N/A"
	if len(languages) > 0 {
		languageDisplay = strings.Join(languages, ", ")
	}

	name := r.Name()

	return models.Country{
		Name:                name,
		Slug:                utils.NameToSlug(name),
		CCA2:                r.CCA2(),
		CCA3:                r.CCA3(),
		Capital:             r.Capital(),
		Region:              r.Region(),
		Subregion:           r.Subregion(),
		Population:          r.Population(),
		FormattedPopulation: utils.FormatPopulation(r.Population()),
		FlagURL:             r.FlagPNG(),
		FlagEmoji:           r.FlagEmoji(),
		CurrencyDisplay:     currencyDisplay,
		LanguageDisplay:     languageDisplay,
		Latitude:            r.Latitude(),
		Longitude:           r.Longitude(),
	}
}

// toListItem converts a RawCountry into a lightweight CountryListItem.
func (s *CountryService) toListItem(r utils.RawCountry) models.CountryListItem {
	c := s.toCountry(r)
	return c.ToListItem()
}