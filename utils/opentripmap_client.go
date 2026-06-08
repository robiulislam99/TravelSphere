// utils/opentripmap_client.go
// OpenTripMapClient is a reusable HTTP client for the OpenTripMap API.
// Docs: https://opentripmap.io/docs
//
// API key is loaded from the OPENTRIPMAP_API_KEY environment variable.
// Transformation into application models.Attraction is done by AttractionService.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const openTripMapBaseURL = "https://api.opentripmap.com/0.1/en"

// OpenTripMapClient wraps an http.Client for the OpenTripMap API.
type OpenTripMapClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewOpenTripMapClient creates a client that reads OPENTRIPMAP_API_KEY from env.
// Returns an error if the API key is not set.
func NewOpenTripMapClient() (*OpenTripMapClient, error) {
	apiKey := os.Getenv("OPENTRIPMAP_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENTRIPMAP_API_KEY environment variable is not set")
	}
	return &OpenTripMapClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    openTripMapBaseURL,
		apiKey:     apiKey,
	}, nil
}

// --- Raw API response types ---

// RawAttractionList is the JSON response from the /places/radius endpoint.
type RawAttractionList struct {
	Type     string             `json:"type"`
	Features []RawAttractionFeature `json:"features"`
}

// RawAttractionFeature represents one GeoJSON feature in the list response.
type RawAttractionFeature struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Properties struct {
		XID  string `json:"xid"`
		Name string `json:"name"`
		Dist float64 `json:"dist"` // distance in meters from search center
		Rate int    `json:"rate"`
		Kinds string `json:"kinds"` // comma-separated category tags
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
	} `json:"geometry"`
}

// RawAttractionDetail is the response from the /places/xid/:xid endpoint.
type RawAttractionDetail struct {
	XID  string `json:"xid"`
	Name string `json:"name"`
	Kinds string `json:"kinds"`
	Point struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"point"`
	Wikipedia string `json:"wikipedia"`
	Image     string `json:"image"`
	Info      struct {
		Descr string `json:"descr"`
	} `json:"info"`
}

// GetAttractionsByRadius fetches attractions within a given radius of lat/lon.
//
// Parameters:
//   lat, lon — coordinates of the country capital or center point
//   radius   — search radius in meters (max 50000 for free tier)
//   limit    — maximum number of results to return
//   kinds    — comma-separated category filter e.g. "interesting_places,historic"
func (c *OpenTripMapClient) GetAttractionsByRadius(lat, lon float64, radius, limit int, kinds string) (*RawAttractionList, error) {
	if kinds == "" {
		kinds = "interesting_places,historic,cultural,museums,architecture,natural"
	}

	params := url.Values{}
	params.Set("radius", fmt.Sprintf("%d", radius))
	params.Set("lon", fmt.Sprintf("%f", lon))
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("kinds", kinds)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("rate", "2")   // minimum rating threshold
	params.Set("format", "geojson")
	params.Set("apikey", c.apiKey)

	reqURL := fmt.Sprintf("%s/places/radius?%s", c.baseURL, params.Encode())

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}

	var result RawAttractionList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode OpenTripMap list response: %w", err)
	}
	return &result, nil
}

// GetAttractionDetail fetches full details for a single attraction by its XID.
func (c *OpenTripMapClient) GetAttractionDetail(xid string) (*RawAttractionDetail, error) {
	if xid == "" {
		return nil, fmt.Errorf("xid must not be empty")
	}
	reqURL := fmt.Sprintf("%s/places/xid/%s?apikey=%s", c.baseURL, xid, c.apiKey)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}

	var detail RawAttractionDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to decode OpenTripMap detail response: %w", err)
	}
	return &detail, nil
}

// doGet is the shared internal HTTP GET helper.
// Returns the raw response body bytes or a descriptive error.
func (c *OpenTripMapClient) doGet(reqURL string) ([]byte, error) {
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("OpenTripMap request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("OpenTripMap API key is invalid or missing")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("OpenTripMap rate limit exceeded; please try again later")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenTripMap API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenTripMap response body: %w", err)
	}
	return body, nil
}