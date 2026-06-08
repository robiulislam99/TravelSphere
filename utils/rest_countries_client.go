// utils/rest_countries_client.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const restCountriesBaseURL = "https://restcountries.com/v3.1"

const defaultFields = "name,cca2,cca3,capital,region,subregion,population,flags,flag,currencies,languages,latlng"

type RestCountriesClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewRestCountriesClient() *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    restCountriesBaseURL,
	}
}

// --- Raw API response types ---

type RawCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA2      string   `json:"cca2"`
	CCA3      string   `json:"cca3"`
	Capital   []string `json:"capital"`
	Region    string   `json:"region"`
	Subregion string   `json:"subregion"`
	Population int64   `json:"population"`
	Flags struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`
	Flag       string `json:"flag"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages map[string]string `json:"languages"`
	Latlng    []float64         `json:"latlng"`
}

// GetAll fetches all countries.
// GET /v3.1/all?fields=...
func (c *RestCountriesClient) GetAll() ([]RawCountry, error) {
	u := fmt.Sprintf("%s/all?fields=%s", c.baseURL, defaultFields)
	return c.fetchCountries(u)
}

// GetByName searches by common or official country name (partial match).
// GET /v3.1/name/{name}
func (c *RestCountriesClient) GetByName(name string) ([]RawCountry, error) {
	if name == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/name/%s?fields=%s", c.baseURL, url.PathEscape(name), defaultFields)
	countries, err := c.fetchCountries(u)
	if err != nil {
		// 404 means no match — not a real error.
		return []RawCountry{}, nil
	}
	return countries, nil
}

// GetByFullName searches by exact full name (common or official).
// GET /v3.1/name/{name}?fullText=true
func (c *RestCountriesClient) GetByFullName(name string) ([]RawCountry, error) {
	if name == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/name/%s?fullText=true&fields=%s", c.baseURL, url.PathEscape(name), defaultFields)
	countries, err := c.fetchCountries(u)
	if err != nil {
		return []RawCountry{}, nil
	}
	return countries, nil
}

// GetByCode searches by a single cca2, cca3, ccn3, or cioc code.
// GET /v3.1/alpha/{code}
func (c *RestCountriesClient) GetByCode(code string) (*RawCountry, error) {
	if code == "" {
		return nil, fmt.Errorf("country code must not be empty")
	}
	u := fmt.Sprintf("%s/alpha/%s?fields=%s", c.baseURL, url.PathEscape(code), defaultFields)
	countries, err := c.fetchCountries(u)
	if err != nil {
		return nil, err
	}
	if len(countries) == 0 {
		return nil, fmt.Errorf("country not found for code: %s", code)
	}
	return &countries[0], nil
}

// GetByCodes searches by multiple country codes at once.
// GET /v3.1/alpha?codes={code},{code}
func (c *RestCountriesClient) GetByCodes(codes []string) ([]RawCountry, error) {
	if len(codes) == 0 {
		return []RawCountry{}, nil
	}
	joined := url.QueryEscape(strings.Join(codes, ","))
	u := fmt.Sprintf("%s/alpha?codes=%s&fields=%s", c.baseURL, joined, defaultFields)
	return c.fetchCountries(u)
}

// GetByCurrency searches by currency code or name.
// GET /v3.1/currency/{currency}
func (c *RestCountriesClient) GetByCurrency(currency string) ([]RawCountry, error) {
	if currency == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/currency/%s?fields=%s", c.baseURL, url.PathEscape(currency), defaultFields)
	return c.fetchCountries(u)
}

// GetByLanguage searches by language code or name.
// GET /v3.1/lang/{language}
func (c *RestCountriesClient) GetByLanguage(language string) ([]RawCountry, error) {
	if language == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/lang/%s?fields=%s", c.baseURL, url.PathEscape(language), defaultFields)
	return c.fetchCountries(u)
}

// GetByCapital searches by capital city name.
// GET /v3.1/capital/{capital}
func (c *RestCountriesClient) GetByCapital(capital string) ([]RawCountry, error) {
	if capital == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/capital/%s?fields=%s", c.baseURL, url.PathEscape(capital), defaultFields)
	return c.fetchCountries(u)
}

// GetByRegion filters countries by region (e.g. "europe", "asia").
// GET /v3.1/region/{region}
func (c *RestCountriesClient) GetByRegion(region string) ([]RawCountry, error) {
	if region == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/region/%s?fields=%s", c.baseURL, url.PathEscape(region), defaultFields)
	return c.fetchCountries(u)
}

// GetBySubregion filters countries by subregion (e.g. "Northern Europe").
// GET /v3.1/subregion/{subregion}
func (c *RestCountriesClient) GetBySubregion(subregion string) ([]RawCountry, error) {
	if subregion == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/subregion/%s?fields=%s", c.baseURL, url.PathEscape(subregion), defaultFields)
	return c.fetchCountries(u)
}

// GetByDemonym searches by demonym (e.g. "peruvian").
// GET /v3.1/demonym/{demonym}
func (c *RestCountriesClient) GetByDemonym(demonym string) ([]RawCountry, error) {
	if demonym == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/demonym/%s?fields=%s", c.baseURL, url.PathEscape(demonym), defaultFields)
	return c.fetchCountries(u)
}

// GetByTranslation searches by any translated country name (e.g. "alemania").
// GET /v3.1/translation/{translation}
func (c *RestCountriesClient) GetByTranslation(translation string) ([]RawCountry, error) {
	if translation == "" {
		return []RawCountry{}, nil
	}
	u := fmt.Sprintf("%s/translation/%s?fields=%s", c.baseURL, url.PathEscape(translation), defaultFields)
	return c.fetchCountries(u)
}

// GetIndependent returns all independent or non-independent countries.
// GET /v3.1/independent?status=true
func (c *RestCountriesClient) GetIndependent(status bool) ([]RawCountry, error) {
	u := fmt.Sprintf("%s/independent?status=%t&fields=%s", c.baseURL, status, defaultFields)
	return c.fetchCountries(u)
}

// fetchCountries is the shared internal helper for all endpoints.
// Every REST Countries endpoint returns a JSON array — including /alpha/:code.
func (c *RestCountriesClient) fetchCountries(u string) ([]RawCountry, error) {
	resp, err := c.httpClient.Get(u)
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