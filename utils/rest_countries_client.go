// utils/rest_countries_client.go

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const restCountriesBaseURL = "https://api.restcountries.com/countries/v5"
const flagsCDNBaseURL = "https://flags.restcountries.com/v5/w320"

// Drop only the heaviest branches — everything else stays nested as documented.
const omitFields = "names.translations,leaders"

const pageSize = 100 // v5 max limit per request

// RawCountry is a single country record as the raw nested JSON object.
type RawCountry map[string]interface{}

// --- Generic nested-path helpers --------------------------------------------

// getNested walks a chain of map keys, e.g. getNested(r, "names", "common").
// Returns nil if any step is missing or not a map.
func getNested(m map[string]interface{}, path ...string) interface{} {
	var cur interface{} = m
	for _, key := range path {
		asMap, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur, ok = asMap[key]
		if !ok {
			return nil
		}
	}
	return cur
}

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func asFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// --- Typed accessors on RawCountry ------------------------------------------

// Name returns the common country name (names.common).
func (r RawCountry) Name() string {
	return asString(getNested(r, "names", "common"))
}

// CCA2 returns the ISO alpha-2 code (codes.alpha_2).
func (r RawCountry) CCA2() string {
	return asString(getNested(r, "codes", "alpha_2"))
}

// CCA3 returns the ISO alpha-3 code (codes.alpha_3).
func (r RawCountry) CCA3() string {
	return asString(getNested(r, "codes", "alpha_3"))
}

// FlagEmoji returns the emoji flag (flag.emoji).
func (r RawCountry) FlagEmoji() string {
	return asString(getNested(r, "flag", "emoji"))
}

// FlagPNG returns a flag image URL from the REST Countries Flags CDN,
// built from the alpha-2 code. The CDN requires no auth and is always available.
// e.g. https://flags.restcountries.com/v5/w320/bd.png
func (r RawCountry) FlagPNG() string {
	code := strings.ToLower(r.CCA2())
	if code == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s.png", flagsCDNBaseURL, code)
}

// Region returns the geographic region (top-level field).
func (r RawCountry) Region() string {
	return asString(r["region"])
}

// Subregion returns the geographic subregion (top-level field).
func (r RawCountry) Subregion() string {
	return asString(r["subregion"])
}

// Population returns the population as int64 (top-level field).
func (r RawCountry) Population() int64 {
	return int64(asFloat(r["population"]))
}

// Latitude returns coordinates.lat.
func (r RawCountry) Latitude() float64 {
	return asFloat(getNested(r, "coordinates", "lat"))
}

// Longitude returns coordinates.lng.
func (r RawCountry) Longitude() float64 {
	return asFloat(getNested(r, "coordinates", "lng"))
}

// Capital returns the primary capital's name from the capitals[] array.
// v5 shape: "capitals": [{"name": "Dhaka", "coordinates": {...}, "attributes": {...}}]
func (r RawCountry) Capital() string {
	caps, ok := r["capitals"].([]interface{})
	if !ok || len(caps) == 0 {
		return ""
	}
	if first, ok := caps[0].(map[string]interface{}); ok {
		return asString(first["name"])
	}
	return ""
}

// Currencies returns a map of currency-code → currency name.
// v5 shape: "currencies": {"BDT": {"name": "Bangladeshi taka", "symbol": "৳"}}
// Falls back to the currency code itself if "name" is missing.
func (r RawCountry) Currencies() map[string]string {
	result := map[string]string{}
	cur, ok := r["currencies"].(map[string]interface{})
	if !ok {
		return result
	}
	for code, v := range cur {
		obj, ok := v.(map[string]interface{})
		if !ok {
			result[code] = code
			continue
		}
		if name := asString(obj["name"]); name != "" {
			result[code] = name
			continue
		}
		result[code] = code
	}
	return result
}

// Languages returns a slice of language display names.
// v5 shape: "languages": [{"name": "Bengali", "native": "বাংলা", "bcp_47": "bn", ...}, ...]
// Tries several candidate keys in case the exact field name differs.
func (r RawCountry) Languages() []string {
	result := []string{}
	langs, ok := r["languages"].([]interface{})
	if !ok {
		return result
	}
	for _, v := range langs {
		obj, ok := v.(map[string]interface{})
		if !ok {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
			continue
		}
		for _, key := range []string{"name", "english_name", "common", "native"} {
			if s := asString(obj[key]); s != "" {
				result = append(result, s)
				break
			}
		}
	}
	return result
}

// --- API response envelope --------------------------------------------------

type apiMeta struct {
	Total  int  `json:"total"`
	Count  int  `json:"count"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	More   bool `json:"more"`
}

type apiResponse struct {
	Data *struct {
		Objects []RawCountry `json:"objects"`
		Meta    apiMeta      `json:"meta"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// --- Client ------------------------------------------------------------------

// RestCountriesClient wraps an http.Client for the REST Countries v5 API.
type RestCountriesClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string

	cacheMu  sync.Mutex
	cache    []RawCountry
	cachedAt time.Time

	debugOnce sync.Once
}

const cacheTTL = 1 * time.Hour

// NewRestCountriesClient creates a client using REST_COUNTRIES_API_KEY from env.
func NewRestCountriesClient() *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		baseURL:    restCountriesBaseURL,
		apiKey:     os.Getenv("REST_COUNTRIES_API_KEY"),
	}
}

// NewRestCountriesClientWithURL — used in tests to point at a mock server.
func NewRestCountriesClientWithURL(baseURL string, apiKey string) *RestCountriesClient {
	if apiKey == "" {
		apiKey = "test-key"
	}
	return &RestCountriesClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// GetAll fetches every country, paginating at 100 per page, with a 1-hour cache.
func (c *RestCountriesClient) GetAll() ([]RawCountry, error) {
	c.cacheMu.Lock()
	if c.cache != nil && time.Since(c.cachedAt) < cacheTTL {
		cached := c.cache
		c.cacheMu.Unlock()
		return cached, nil
	}
	c.cacheMu.Unlock()

	var all []RawCountry
	offset := 0

	for {
		reqURL := fmt.Sprintf("%s?response_fields_omit=%s&limit=%d&offset=%d",
			c.baseURL, url.QueryEscape(omitFields), pageSize, offset)

		page, meta, err := c.fetchPage(reqURL)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)

		if !meta.More {
			break
		}
		offset += pageSize
	}

	// One-time debug: if the first country's name didn't parse, dump its keys
	// so you can see the actual shape returned by the API.
	c.debugOnce.Do(func() {
		if len(all) > 0 && all[0].Name() == "" {
			raw, _ := json.MarshalIndent(map[string]interface{}(all[0]), "", "  ")
			log.Printf("[RestCountriesClient] DEBUG — Name() empty. First object keys/shape:\n%s", truncate(string(raw), 1500))
		}
	})

	c.cacheMu.Lock()
	c.cache = all
	c.cachedAt = time.Now()
	c.cacheMu.Unlock()

	return all, nil
}

// GetByName searches countries whose common name contains the query (substring match).
func (c *RestCountriesClient) GetByName(name string) ([]RawCountry, error) {
	if name == "" {
		return []RawCountry{}, nil
	}
	reqURL := fmt.Sprintf("%s/names.common?q=%s&response_fields_omit=%s&limit=%d",
		c.baseURL, url.QueryEscape(name), url.QueryEscape(omitFields), pageSize)

	countries, _, err := c.fetchPage(reqURL)
	if err != nil {
		return []RawCountry{}, nil // treat errors as "no match" for search UX
	}
	return countries, nil
}

// GetByCode fetches a single country by alpha-2 or alpha-3 code (exact match).
func (c *RestCountriesClient) GetByCode(code string) (*RawCountry, error) {
	if code == "" {
		return nil, fmt.Errorf("country code must not be empty")
	}

	property := "codes.alpha_2"
	if len(code) == 3 {
		property = "codes.alpha_3"
	}

	reqURL := fmt.Sprintf("%s/%s/%s?response_fields_omit=%s",
		c.baseURL, property, url.PathEscape(code), url.QueryEscape(omitFields))

	countries, _, err := c.fetchPage(reqURL)
	if err != nil {
		return nil, err
	}
	if len(countries) == 0 {
		return nil, fmt.Errorf("country not found for code: %s", code)
	}
	return &countries[0], nil
}

// fetchPage performs the GET request and decodes the v5 envelope.
func (c *RestCountriesClient) fetchPage(reqURL string) ([]RawCountry, apiMeta, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, apiMeta{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("User-Agent", "TravelSphere/1.0")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, apiMeta{}, fmt.Errorf("REST Countries request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiMeta{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return []RawCountry{}, apiMeta{}, nil
	}

	var parsed apiResponse
	if jsonErr := json.Unmarshal(body, &parsed); jsonErr != nil {
		return nil, apiMeta{}, fmt.Errorf("failed to decode response (status %d): %s",
			resp.StatusCode, truncate(string(body), 300))
	}

	if len(parsed.Errors) > 0 {
		return nil, apiMeta{}, fmt.Errorf("REST Countries API error (status %d): %s",
			resp.StatusCode, parsed.Errors[0].Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apiMeta{}, fmt.Errorf("REST Countries API returned status %d", resp.StatusCode)
	}

	if parsed.Data == nil {
		return []RawCountry{}, apiMeta{}, nil
	}

	return parsed.Data.Objects, parsed.Data.Meta, nil
}

// truncate limits a string length for safe error message display.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}