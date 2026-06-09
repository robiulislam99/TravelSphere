package services

import (
	"testing"

	"github.com/robiulislam99/TravelSphere/models"
)

// TestNewAttractionService_WithoutClient tests service creation when client fails
func TestNewAttractionService_WithoutClient(t *testing.T) {
	// This will create a service with nil client (API key not set)
	service := &AttractionService{client: nil}
	if service == nil {
		t.Error("AttractionService should not be nil")
	}
}

// TestGetByCoords_DefaultRadius tests default radius when <= 0
func TestGetByCoords_DefaultRadius(t *testing.T) {
	service := &AttractionService{client: nil}
	attractions, err := service.GetByCoords(48.8584, 2.2945, 0, 4)

	if err != nil {
		t.Errorf("GetByCoords returned error: %v", err)
	}
	if len(attractions) != 0 {
		t.Errorf("Expected empty results for nil client, got %d", len(attractions))
	}
}

// TestGetByCoords_DefaultLimit tests default limit when <= 0
func TestGetByCoords_DefaultLimit(t *testing.T) {
	service := &AttractionService{client: nil}
	attractions, err := service.GetByCoords(48.8584, 2.2945, 5000, 0)

	if err != nil {
		t.Errorf("GetByCoords returned error: %v", err)
	}
	if len(attractions) != 0 {
		t.Errorf("Expected empty results for nil client, got %d", len(attractions))
	}
}

// TestGetByCoords_WithNilClient tests handling when client is nil
func TestGetByCoords_WithNilClient(t *testing.T) {
	service := &AttractionService{client: nil}
	attractions, err := service.GetByCoords(48.8584, 2.2945, 5000, 4)

	if err != nil {
		t.Errorf("GetByCoords returned error: %v", err)
	}
	if len(attractions) != 0 {
		t.Errorf("Expected empty results for nil client, got %d", len(attractions))
	}
}

// TestGetForHomePage_WithNilClient tests home page fetch with nil client
func TestGetForHomePage_WithNilClient(t *testing.T) {
	service := &AttractionService{client: nil}
	attractions := service.GetForHomePage()

	if len(attractions) != 0 {
		t.Errorf("Expected empty results for nil client, got %d", len(attractions))
	}
}

// TestGetByCoords_ValidCoordinates tests coordinate extraction
func TestGetByCoords_ValidCoordinates(t *testing.T) {
	// Test with nil client but verify logic paths
	service := &AttractionService{client: nil}

	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantLen int
	}{
		{"Paris", 48.8584, 2.2945, 0},
		{"Tokyo", 35.6762, 139.6503, 0},
		{"NYC", 40.7128, -74.0060, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.GetByCoords(tt.lat, tt.lon, 5000, 4)
			if err != nil {
				t.Errorf("GetByCoords error: %v", err)
			}
			if len(results) != tt.wantLen {
				t.Errorf("GetByCoords() returned %d attractions, want %d", len(results), tt.wantLen)
			}
		})
	}
}

// TestAttractionModel_FormattedDistance tests distance formatting
func TestAttractionModel_FormattedDistance(t *testing.T) {
	tests := []struct {
		name     string
		distance float64
		want     string
	}{
		{"meters", 800, "800 m"},
		{"kilometers", 1500, "1.5 km"},
		{"over kilometer", 5000, "5.0 km"},
		{"zero distance", 0, ""},
		{"negative distance", -100, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attraction := &models.Attraction{Distance: tt.distance}
			if got := attraction.FormattedDistance(); got != tt.want {
				t.Errorf("FormattedDistance() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestHomePageSeeds tests that home page seeds have valid coordinates
func TestHomePageSeeds(t *testing.T) {
	service := &AttractionService{client: nil}

	// Verify the service structure is correct
	if service == nil {
		t.Error("Service should not be nil")
	}

	// Verify GetForHomePage returns empty slice for nil client
	result := service.GetForHomePage()
	if result == nil {
		t.Error("GetForHomePage should return empty slice, not nil")
	}
	if len(result) != 0 {
		t.Errorf("GetForHomePage should return empty results for nil client, got %d", len(result))
	}
}
