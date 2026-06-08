// utils/opentripmap_client.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const openTripMapBaseURL = "https://api.opentripmap.com/0.1/en"

type OpenTripMapClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

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

type RawAttractionList struct {
	Type     string                 `json:"type"`
	Features []RawAttractionFeature `json:"features"`
}

type RawAttractionFeature struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Properties struct {
		XID   string  `json:"xid"`
		Name  string  `json:"name"`
		Dist  float64 `json:"dist"`
		Rate  int     `json:"rate"`
		Kinds string  `json:"kinds"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"` // GeoJSON: [longitude, latitude]
	} `json:"geometry"`
}

type RawAttractionDetail struct {
	XID   string `json:"xid"`
	Name  string `json:"name"`
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

// SafeCoords extracts (lon, lat) from a GeoJSON coordinates slice.
// Exported so AttractionService can use it directly instead of
// duplicating the bounds check.
// Returns (0, 0) if the slice has fewer than 2 elements.
func SafeCoords(coords []float64) (lon, lat float64) {
	if len(coords) < 2 {
		return 0, 0
	}
	return coords[0], coords[1] // GeoJSON order: [longitude, latitude]
}

func (c *OpenTripMapClient) GetAttractionsByRadius(lat, lon float64, radius, limit int, kinds string) (*RawAttractionList, error) {
	if kinds == "" {
		kinds = "interesting_places,historic,cultural,museums,architecture,natural"
	}

	params := url.Values{}
	params.Set("radius", strconv.Itoa(radius))
	params.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64)) // full precision
	params.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	params.Set("kinds", kinds)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("rate", "2")
	params.Set("format", "geojson")
	params.Set("apikey", c.apiKey)

	reqURL := fmt.Sprintf("%s/places/radius?%s", c.baseURL, params.Encode())

	body, err := c.doGetWithRetry(reqURL)
	if err != nil {
		return nil, err
	}

	var result RawAttractionList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode OpenTripMap list response: %w", err)
	}
	return &result, nil
}

func (c *OpenTripMapClient) GetAttractionDetail(xid string) (*RawAttractionDetail, error) {
	if xid == "" {
		return nil, fmt.Errorf("xid must not be empty")
	}

	// apikey goes through params.Encode() so it doesn't appear raw in a logged URL string
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	reqURL := fmt.Sprintf("%s/places/xid/%s?%s", c.baseURL, xid, params.Encode())

	body, err := c.doGetWithRetry(reqURL)
	if err != nil {
		return nil, err
	}

	var detail RawAttractionDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to decode OpenTripMap detail response: %w", err)
	}
	return &detail, nil
}

// doGetWithRetry attempts the request up to 2 times on 5xx errors.
func (c *OpenTripMapClient) doGetWithRetry(reqURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		body, err := c.doGet(reqURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		// Only retry on server-side errors, not on 401/429/bad input
		if isRetryable(err) && attempt == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	return nil, lastErr
}

// doGet performs a single HTTP GET and returns the body bytes.
func (c *OpenTripMapClient) doGet(reqURL string) ([]byte, error) {
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("OpenTripMap request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("OpenTripMap API key is invalid or missing")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("OpenTripMap rate limit exceeded; please try again later")
	case http.StatusOK:
		// continue below
	default:
		return nil, fmt.Errorf("OpenTripMap API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenTripMap response body: %w", err)
	}
	return body, nil
}

// isRetryable returns true for transient server-side errors worth retrying.
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 0 && (contains(msg, "status 5") || contains(msg, "request failed"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
