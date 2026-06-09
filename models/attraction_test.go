package models

import (
	"testing"
)

func TestAttractionFormattedDistance(t *testing.T) {
	tests := []struct {
		name     string
		distance float64
		expected string
	}{
		{
			name:     "zero distance returns empty string",
			distance: 0,
			expected: "",
		},
		{
			name:     "negative distance returns empty string",
			distance: -500,
			expected: "",
		},
		{
			name:     "distance less than 1000 meters",
			distance: 500,
			expected: "500 m",
		},
		{
			name:     "distance exactly 1000 meters",
			distance: 1000,
			expected: "1.0 km",
		},
		{
			name:     "distance in kilometers",
			distance: 1500,
			expected: "1.5 km",
		},
		{
			name:     "large distance in kilometers",
			distance: 5250,
			expected: "5.2 km",
		},
		{
			name:     "very small distance in meters",
			distance: 50,
			expected: "50 m",
		},
		{
			name:     "distance with decimal point less than 1000",
			distance: 750.5,
			expected: "750 m", // rounded down due to %.0f format
		},
		{
			name:     "distance with decimal point greater than 1000",
			distance: 1500.5,
			expected: "1.5 km",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Attraction{Distance: tt.distance}
			result := a.FormattedDistance()
			if result != tt.expected {
				t.Errorf("FormattedDistance() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestAttractionNewInstance(t *testing.T) {
	a := &Attraction{
		XID:         "xid123",
		Name:        "Eiffel Tower",
		Kinds:       "historic,architecture",
		PrimaryKind: "historic",
		Latitude:    48.8584,
		Longitude:   2.2945,
		Distance:    2500,
		CountryName: "France",
		Description: "Famous iron tower",
		WikipediaURL: "https://en.wikipedia.org/wiki/Eiffel_Tower",
		ImageURL:    "https://example.com/eiffel.jpg",
	}

	if a.XID != "xid123" {
		t.Errorf("XID = %s, want xid123", a.XID)
	}
	if a.Name != "Eiffel Tower" {
		t.Errorf("Name = %s, want Eiffel Tower", a.Name)
	}
	if a.Distance != 2500 {
		t.Errorf("Distance = %f, want 2500", a.Distance)
	}
	if a.CountryName != "France" {
		t.Errorf("CountryName = %s, want France", a.CountryName)
	}
}
