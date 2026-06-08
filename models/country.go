// models/country.go
// Country represents a single country entity returned by the REST Countries API
// and used throughout the application (controllers, services, templates).
package models

// Country holds all display-ready fields for a country.
// Raw API data is transformed into this struct by the CountryService.
type Country struct {
	// Core identity
	Name    string `json:"name"`
	Slug    string `json:"slug"`    // lowercase, hyphen-separated e.g. "united-states"
	CCA2    string `json:"cca2"`    // ISO 3166-1 alpha-2 code e.g. "US"
	CCA3    string `json:"cca3"`    // ISO 3166-1 alpha-3 code e.g. "USA"
	Region  string `json:"region"`  // e.g. "Asia", "Europe"
	Subregion string `json:"subregion"`

	// Display fields
	Capital             string  `json:"capital"`
	Population          int64   `json:"population"`
	FormattedPopulation string  `json:"formatted_population"` // e.g. "1.4B"
	FlagURL             string  `json:"flag_url"`             // PNG URL from REST Countries
	FlagEmoji           string  `json:"flag_emoji"`

	// Currency and language as display strings
	CurrencyDisplay string `json:"currency_display"` // e.g. "BDT (Bangladeshi Taka)"
	LanguageDisplay string `json:"language_display"` // e.g. "Bengali, English"

	// Raw data kept for service-layer use (lat/lon for OpenTripMap)
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CountryListItem is a lighter version used in list/grid views and AJAX responses.
type CountryListItem struct {
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Capital             string `json:"capital"`
	Region              string `json:"region"`
	FlagURL             string `json:"flag_url"`
	FlagEmoji           string `json:"flag_emoji"`
	FormattedPopulation string `json:"formatted_population"`
	CurrencyDisplay     string `json:"currency_display"`
	LanguageDisplay     string `json:"language_display"`
}

// ToListItem converts a full Country into a CountryListItem for list views.
func (c *Country) ToListItem() CountryListItem {
	return CountryListItem{
		Name:                c.Name,
		Slug:                c.Slug,
		Capital:             c.Capital,
		Region:              c.Region,
		FlagURL:             c.FlagURL,
		FlagEmoji:           c.FlagEmoji,
		FormattedPopulation: c.FormattedPopulation,
		CurrencyDisplay:     c.CurrencyDisplay,
		LanguageDisplay:     c.LanguageDisplay,
	}
}