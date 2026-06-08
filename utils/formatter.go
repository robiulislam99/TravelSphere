// utils/formatter.go
// Formatting helpers that convert raw API data into human-readable strings.
// Used by CountryService when building models.Country structs.
package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// FormatPopulation converts a raw int64 population count into a short readable string.
// Examples:
//   1_400_000_000 → "1.4B"
//   170_000_000   → "170M"
//   500_000       → "500K"
//   12_345        → "12,345"
func FormatPopulation(pop int64) string {
	switch {
	case pop >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(pop)/1_000_000_000)
	case pop >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(pop)/1_000_000)
	case pop >= 1_000:
		return fmt.Sprintf("%.1fK", float64(pop)/1_000)
	default:
		return fmt.Sprintf("%d", pop)
	}
}

// FormatCurrencies converts the REST Countries currency map into a display string.
// Input: map[code]{ name, symbol } e.g. {"BDT": {name:"Bangladeshi Taka", symbol:"৳"}}
// Output: "BDT (Bangladeshi Taka)" — or multiple joined by ", "
func FormatCurrencies(currencies map[string]struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}) string {
	if len(currencies) == 0 {
		return "N/A"
	}
	parts := make([]string, 0, len(currencies))
	for code, cur := range currencies {
		if cur.Name != "" {
			parts = append(parts, fmt.Sprintf("%s (%s)", code, cur.Name))
		} else {
			parts = append(parts, code)
		}
	}
	return strings.Join(parts, ", ")
}

// FormatLanguages converts the REST Countries language map into a display string.
// Input: map[code]name e.g. {"ben": "Bengali", "eng": "English"}
// Output: "Bengali, English"
func FormatLanguages(languages map[string]string) string {
	if len(languages) == 0 {
		return "N/A"
	}
	names := make([]string, 0, len(languages))
	for _, name := range languages {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// NameToSlug converts a country name to a URL-safe slug.
// Examples:
//   "United States"     → "united-states"
//   "Côte d'Ivoire"     → "cote-d-ivoire"
//   "São Tomé & Príncipe" → "sao-tome-principe"
func NameToSlug(name string) string {
	// Normalize: lowercase
	slug := strings.ToLower(name)

	// Replace common accented characters
	replacements := map[string]string{
		"à": "a", "á": "a", "â": "a", "ã": "a", "ä": "a", "å": "a",
		"è": "e", "é": "e", "ê": "e", "ë": "e",
		"ì": "i", "í": "i", "î": "i", "ï": "i",
		"ò": "o", "ó": "o", "ô": "o", "õ": "o", "ö": "o",
		"ù": "u", "ú": "u", "û": "u", "ü": "u",
		"ñ": "n", "ç": "c", "ý": "y",
		"ā": "a", "ē": "e", "ī": "i", "ō": "o", "ū": "u",
	}
	for accented, plain := range replacements {
		slug = strings.ReplaceAll(slug, accented, plain)
	}

	// Replace non-alphanumeric characters (spaces, punctuation) with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// PrimaryKind extracts the first meaningful category from a comma-separated kinds string.
// e.g. "historic,architecture,fortifications" → "historic"
func PrimaryKind(kinds string) string {
	if kinds == "" {
		return "attraction"
	}
	parts := strings.Split(kinds, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		// Skip overly generic tags
		if p != "" && p != "interesting_places" && p != "other" {
			return strings.ReplaceAll(p, "_", " ")
		}
	}
	return strings.ReplaceAll(strings.TrimSpace(parts[0]), "_", " ")
}

// TruncateString returns the first n characters of s, appending "…" if truncated.
func TruncateString(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}