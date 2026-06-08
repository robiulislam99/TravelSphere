// utils/rest_countries_client.go
// RestCountriesClient is a reusable HTTP client for the REST Countries v3.1 API.
// Docs: https://restcountries.com/
//
// All methods return raw decoded API responses; transformation into
// application models is done by CountryService (services/country_service.go).
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const restCountriesBaseURL = "https://restcountries.com/v3.1"

// RestCountriesClient wraps an http.Client for the REST Countries API.
type RestCountriesClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewRestCountriesClient creates a client with a sensible default timeout.
func NewRestCountriesClient() *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    restCountriesBaseURL,
	}
}

// --- Raw API response types ---
// These mirror the REST Countries v3.1 JSON structure exactly.
// CountryService maps these into application models.Country structs.

// RawCountry is the raw JSON shape returned by the REST Countries API.
type RawCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA2    string   `json:"cca2"`
	CCA3    string   `json:"cca3"`
	Capital []string `json:"capital"`
	Region  string   `json:"region"`
	Subregion string `json:"subregion"`
	Population int64  `json:"population"`
	Flags   struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`
	Flag      string                        `json:"flag"` // emoji
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages  map[string]string             `json:"languages"`
	Latlng     []float64                     `json:"latlng"`
}

// GetAll fetches all countries from the REST Countries API.
// Returns a slice of raw country data and any error encountered.
func (c *RestCountriesClient) GetAll() ([]RawCountry, error) {
	// Request only the fields we need to keep the response lean
	fields := "name,cca2,cca3,capital,region,subregion,population,flags,flag,currencies,languages,latlng"
	url := fmt.Sprintf("%s/all?fields=%s", c.baseURL, fields)
	return c.fetchCountries(url)
}

// GetByName fetches countries matching a name query (partial match supported).
// Returns an empty slice (not an error) when no countries match.
func (c *RestCountriesClient) GetByName(name string) ([]RawCountry, error) {
	if name == "" {
		return []RawCountry{}, nil
	}
	fields := "name,cca2,cca3,capital,region,subregion,population,flags,flag,currencies,languages,latlng"
	url := fmt.Sprintf("%s/name/%s?fields=%s", c.baseURL, name, fields)

	countries, err := c.fetchCountries(url)
	if err != nil {
		// REST Countries returns 404 when no match — treat as empty result
		return []RawCountry{}, nil
	}
	return countries, nil
}

// GetByCode fetches a single country by its alpha-2 or alpha-3 code.
// e.g. "BD" or "BGD" for Bangladesh.
func (c *RestCountriesClient) GetByCode(code string) (*RawCountry, error) {
	if code == "" {
		return nil, fmt.Errorf("country code must not be empty")
	}
	fields := "name,cca2,cca3,capital,region,subregion,population,flags,flag,currencies,languages,latlng"
	url := fmt.Sprintf("%s/alpha/%s?fields=%s", c.baseURL, code, fields)

	countries, err := c.fetchCountries(url)
	if err != nil {
		return nil, err
	}
	if len(countries) == 0 {
		return nil, fmt.Errorf("country not found for code: %s", code)
	}
	return &countries[0], nil
}

// fetchCountries is the shared internal helper that performs the HTTP GET,
// reads the body, and decodes the JSON array of countries.
func (c *RestCountriesClient) fetchCountries(url string) ([]RawCountry, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("rest countries request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []RawCountry{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rest countries API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read REST countries response body: %w", err)
	}

	var countries []RawCountry
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to decode REST countries response: %w", err)
	}

	return countries, nil
}