// utils/rest_countries_client.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const restCountriesBaseURL = "https://restcountries.com/v3.1"

// REST Countries limits field selections to 10 fields for this endpoint.
const restCountriesFields = "name,cca2,capital,region,subregion,population,flags,currencies,languages,latlng"

type RestCountriesClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewRestCountriesClient() *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    restCountriesBaseURL,
	}
}

// NewRestCountriesClientWithURL creates a client pointed at a custom base URL (used in tests).
func NewRestCountriesClientWithURL(baseURL string) *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

// RawCountry mirrors the REST Countries v3.1 JSON shape exactly.
type RawCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA2       string   `json:"cca2"`
	CCA3       string   `json:"cca3"`
	Capital    []string `json:"capital"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Population int64    `json:"population"`
	Flags      struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`
	Flag       string `json:"flag"` // emoji
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages map[string]string `json:"languages"`
	Latlng    []float64         `json:"latlng"`
}

// GetAll fetches every country with only the fields needed by the app.
func (c *RestCountriesClient) GetAll() ([]RawCountry, error) {
	url := fmt.Sprintf("%s/all?fields=%s", c.baseURL, restCountriesFields)
	return c.fetchCountries(url)
}

// GetByName fetches countries whose name matches the query.
func (c *RestCountriesClient) GetByName(name string) ([]RawCountry, error) {
	if name == "" {
		return []RawCountry{}, nil
	}
	url := fmt.Sprintf("%s/name/%s?fields=%s", c.baseURL, url.PathEscape(name), restCountriesFields)
	countries, err := c.fetchCountries(url)
	if err != nil {
		return []RawCountry{}, nil // treat not-found as empty
	}
	return countries, nil
}

// GetByCode fetches a single country by alpha-2 or alpha-3 code.
func (c *RestCountriesClient) GetByCode(code string) (*RawCountry, error) {
	if code == "" {
		return nil, fmt.Errorf("country code must not be empty")
	}
	url := fmt.Sprintf("%s/alpha/%s?fields=%s", c.baseURL, url.PathEscape(code), restCountriesFields)
	countries, err := c.fetchCountries(url)
	if err != nil {
		return nil, err
	}
	if len(countries) == 0 {
		return nil, fmt.Errorf("country not found for code: %s", code)
	}
	return &countries[0], nil
}

// fetchCountries is the shared internal GET + decode helper.
func (c *RestCountriesClient) fetchCountries(url string) ([]RawCountry, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	// Some CDN edges block requests without a User-Agent
	req.Header.Set("User-Agent", "TravelSphere/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("REST Countries request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []RawCountry{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("REST Countries API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var countries []RawCountry
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return countries, nil
}
