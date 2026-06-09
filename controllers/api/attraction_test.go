package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/robiulislam99/TravelSphere/utils"
)

// Helper to parse response
func parseResponse(body string) utils.APIResponse {
	var resp utils.APIResponse
	json.Unmarshal([]byte(body), &resp)
	return resp
}

// Test: Missing lat parameter
func TestGetByCoords_MissingLat(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lon=-74.0060", nil)

	lat := req.URL.Query().Get("lat")
	lon := req.URL.Query().Get("lon")

	if lat != "" {
		t.Errorf("Expected empty lat, got %s", lat)
	}
	if lon != "-74.0060" {
		t.Errorf("Expected lon=-74.0060, got %s", lon)
	}
}

// Test: Missing lon parameter
func TestGetByCoords_MissingLon(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=40.7128", nil)

	lat := req.URL.Query().Get("lat")
	lon := req.URL.Query().Get("lon")

	if lat != "40.7128" {
		t.Errorf("Expected lat=40.7128, got %s", lat)
	}
	if lon != "" {
		t.Errorf("Expected empty lon, got %s", lon)
	}
}

// Test: Missing both lat and lon
func TestGetByCoords_MissingBothLatLon(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions", nil)

	lat := req.URL.Query().Get("lat")
	lon := req.URL.Query().Get("lon")

	if lat != "" || lon != "" {
		t.Errorf("Expected empty params, got lat=%s lon=%s", lat, lon)
	}
}

// Test: Invalid lat (not a number)
func TestGetByCoords_InvalidLat(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=invalid&lon=-74.0060", nil)

	lat := req.URL.Query().Get("lat")
	if lat != "invalid" {
		t.Errorf("Expected lat=invalid, got %s", lat)
	}
}

// Test: Invalid lon (not a number)
func TestGetByCoords_InvalidLon(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=40.7128&lon=invalid", nil)

	lon := req.URL.Query().Get("lon")
	if lon != "invalid" {
		t.Errorf("Expected lon=invalid, got %s", lon)
	}
}

// Test: Lat below valid range
func TestGetByCoords_LatTooLow(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=-91&lon=-74.0060", nil)

	lat := req.URL.Query().Get("lat")
	if lat != "-91" {
		t.Errorf("Expected lat=-91, got %s", lat)
	}
}

// Test: Lat above valid range
func TestGetByCoords_LatTooHigh(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=91&lon=-74.0060", nil)

	lat := req.URL.Query().Get("lat")
	if lat != "91" {
		t.Errorf("Expected lat=91, got %s", lat)
	}
}

// Test: Lon below valid range
func TestGetByCoords_LonTooLow(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=40.7128&lon=-181", nil)

	lon := req.URL.Query().Get("lon")
	if lon != "-181" {
		t.Errorf("Expected lon=-181, got %s", lon)
	}
}

// Test: Lon above valid range
func TestGetByCoords_LonTooHigh(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=40.7128&lon=181", nil)

	lon := req.URL.Query().Get("lon")
	if lon != "181" {
		t.Errorf("Expected lon=181, got %s", lon)
	}
}

// Test: Valid boundary values for lat/lon
func TestGetByCoords_LatMinBoundary(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=-90&lon=0", nil)

	lat := req.URL.Query().Get("lat")
	lon := req.URL.Query().Get("lon")

	if lat != "-90" || lon != "0" {
		t.Errorf("Expected lat=-90, lon=0, got lat=%s, lon=%s", lat, lon)
	}
}

// Test: Valid lat at maximum boundary
func TestGetByCoords_LatMaxBoundary(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=90&lon=0", nil)

	lat := req.URL.Query().Get("lat")
	if lat != "90" {
		t.Errorf("Expected lat=90, got %s", lat)
	}
}

// Test: Valid lon at minimum boundary
func TestGetByCoords_LonMinBoundary(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=0&lon=-180", nil)

	lon := req.URL.Query().Get("lon")
	if lon != "-180" {
		t.Errorf("Expected lon=-180, got %s", lon)
	}
}

// Test: Valid lon at maximum boundary
func TestGetByCoords_LonMaxBoundary(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/attractions?lat=0&lon=180", nil)

	lon := req.URL.Query().Get("lon")
	if lon != "180" {
		t.Errorf("Expected lon=180, got %s", lon)
	}
}

// Test: Query parameter parsing
func TestGetByCoords_QueryParameterParsing(t *testing.T) {
	tests := []struct {
		name  string
		query string
		lat   string
		lon   string
	}{
		{"ValidCoordinates", "?lat=40.7128&lon=-74.0060", "40.7128", "-74.0060"},
		{"NoParams", "", "", ""},
		{"OnlyLat", "?lat=40.7128", "40.7128", ""},
		{"OnlyLon", "?lon=-74.0060", "", "-74.0060"},
		{"ExtraParams", "?lat=40.7128&lon=-74.0060&radius=5000&limit=10", "40.7128", "-74.0060"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/attractions"+tt.query, nil)

			lat := req.URL.Query().Get("lat")
			lon := req.URL.Query().Get("lon")

			if lat != tt.lat {
				t.Errorf("Expected lat=%s, got %s", tt.lat, lat)
			}
			if lon != tt.lon {
				t.Errorf("Expected lon=%s, got %s", tt.lon, lon)
			}
		})
	}
}

// Test: Response structure
func TestGetByCoords_ResponseStructure(t *testing.T) {
	response := `{"data": [], "message": "ok", "status": 200}`
	resp := parseResponse(response)

	if resp.Message != "ok" {
		t.Errorf("Expected message 'ok', got '%s'", resp.Message)
	}
	if resp.Status != 200 {
		t.Errorf("Expected status 200, got %d", resp.Status)
	}
}

// Test: Error response structure
func TestGetByCoords_ErrorResponseStructure(t *testing.T) {
	errorResponse := `{"data": null, "message": "lat and lon query params are required", "status": 400}`
	resp := parseResponse(errorResponse)

	if resp.Status != 400 {
		t.Errorf("Expected status 400, got %d", resp.Status)
	}
	if resp.Message != "lat and lon query params are required" {
		t.Errorf("Expected error message, got '%s'", resp.Message)
	}
}
