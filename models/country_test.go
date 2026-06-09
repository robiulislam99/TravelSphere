package models

import (
	"testing"
)

func TestCountryToListItem(t *testing.T) {
	country := &Country{
		Name:                "France",
		Slug:                "france",
		Capital:             "Paris",
		Region:              "Europe",
		FlagURL:             "https://example.com/france.png",
		FlagEmoji:           "🇫🇷",
		FormattedPopulation: "67.4M",
		CurrencyDisplay:     "EUR (Euro)",
		LanguageDisplay:     "French",
		CCA2:                "FR",
		CCA3:                "FRA",
		Subregion:           "Western Europe",
		Population:          67750000,
		Latitude:            46.227638,
		Longitude:           2.213749,
	}

	listItem := country.ToListItem()

	if listItem.Name != "France" {
		t.Errorf("ToListItem().Name = %s, want France", listItem.Name)
	}
	if listItem.Slug != "france" {
		t.Errorf("ToListItem().Slug = %s, want france", listItem.Slug)
	}
	if listItem.Capital != "Paris" {
		t.Errorf("ToListItem().Capital = %s, want Paris", listItem.Capital)
	}
	if listItem.Region != "Europe" {
		t.Errorf("ToListItem().Region = %s, want Europe", listItem.Region)
	}
	if listItem.FlagURL != "https://example.com/france.png" {
		t.Errorf("ToListItem().FlagURL = %s, want https://example.com/france.png", listItem.FlagURL)
	}
	if listItem.FlagEmoji != "🇫🇷" {
		t.Errorf("ToListItem().FlagEmoji = %s, want 🇫🇷", listItem.FlagEmoji)
	}
	if listItem.FormattedPopulation != "67.4M" {
		t.Errorf("ToListItem().FormattedPopulation = %s, want 67.4M", listItem.FormattedPopulation)
	}
	if listItem.CurrencyDisplay != "EUR (Euro)" {
		t.Errorf("ToListItem().CurrencyDisplay = %s, want EUR (Euro)", listItem.CurrencyDisplay)
	}
	if listItem.LanguageDisplay != "French" {
		t.Errorf("ToListItem().LanguageDisplay = %s, want French", listItem.LanguageDisplay)
	}
}

func TestCountryToListItemWithEmptyFields(t *testing.T) {
	country := &Country{
		Name: "TestCountry",
		Slug: "test-country",
	}

	listItem := country.ToListItem()

	if listItem.Name != "TestCountry" {
		t.Errorf("ToListItem().Name = %s, want TestCountry", listItem.Name)
	}
	if listItem.Slug != "test-country" {
		t.Errorf("ToListItem().Slug = %s, want test-country", listItem.Slug)
	}
	if listItem.Capital != "" {
		t.Errorf("ToListItem().Capital = %s, want empty string", listItem.Capital)
	}
	if listItem.FlagURL != "" {
		t.Errorf("ToListItem().FlagURL = %s, want empty string", listItem.FlagURL)
	}
}

func TestCountryNewInstance(t *testing.T) {
	country := &Country{
		Name:                "Japan",
		Slug:                "japan",
		CCA2:                "JP",
		CCA3:                "JPN",
		Region:              "Asia",
		Subregion:           "East Asia",
		Capital:             "Tokyo",
		Population:          125124000,
		FormattedPopulation: "125.1M",
		FlagURL:             "https://example.com/japan.png",
		FlagEmoji:           "🇯🇵",
		CurrencyDisplay:     "JPY (Japanese Yen)",
		LanguageDisplay:     "Japanese",
		Latitude:            36.204823,
		Longitude:           138.252924,
	}

	if country.CCA2 != "JP" {
		t.Errorf("CCA2 = %s, want JP", country.CCA2)
	}
	if country.Population != 125124000 {
		t.Errorf("Population = %d, want 125124000", country.Population)
	}
	if country.Region != "Asia" {
		t.Errorf("Region = %s, want Asia", country.Region)
	}
}
