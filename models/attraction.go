// models/attraction.go
// Attraction represents a tourist attraction, landmark, museum, or historical site
// returned by the OpenTripMap API and transformed for display.
package models

import "fmt"

// Attraction holds display-ready data for a single point of interest.
type Attraction struct {
	// OpenTripMap unique identifier
	XID  string `json:"xid"`
	Name string `json:"name"`

	// Category / type — comma-separated raw kinds from OpenTripMap
	// e.g. "historic,architecture,fortifications"
	Kinds string `json:"kinds"`

	// Primary kind derived from Kinds for display
	PrimaryKind string `json:"primary_kind"`

	// Geographic data
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	// Distance from search center in meters (set by service layer)
	Distance float64 `json:"distance"`

	// Optional: populated when fetching full place details
	CountryName string `json:"country_name"`
	Description string `json:"description"`
	WikipediaURL string `json:"wikipedia_url"`
	ImageURL    string `json:"image_url"`
}

// FormattedDistance returns a human-readable distance string.
// e.g. 1500 → "1.5 km", 800 → "800 m"
func (a *Attraction) FormattedDistance() string {
	if a.Distance <= 0 {
		return ""
	}
	if a.Distance >= 1000 {
		return fmt.Sprintf("%.1f km", a.Distance/1000)
	}
	return fmt.Sprintf("%.0f m", a.Distance)
}