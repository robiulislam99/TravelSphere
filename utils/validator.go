// utils/validator.go
// Validation helpers used across the application.
// Called by service methods and API controllers before processing requests.
package utils

import (
	"strings"
	"unicode"
)

// IsValidSlug returns true if the slug contains only lowercase letters,
// digits, and hyphens — and is not empty.
// e.g. "bangladesh", "united-states", "cote-d-ivoire"
func IsValidSlug(slug string) bool {
	if strings.TrimSpace(slug) == "" {
		return false
	}
	for _, r := range slug {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

// IsValidStatus returns true if the status string is an allowed wishlist status.
// Accepted values: "Planned", "Visited"
func IsValidStatus(status string) bool {
	return status == "Planned" || status == "Visited"
}

// IsNonEmpty returns true when the string is not empty after trimming whitespace.
func IsNonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// SanitizeString trims leading/trailing whitespace from a string.
// Use for cleaning user-provided input before storing or comparing.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// IsValidRegion returns true if the region matches one of the known REST Countries regions.
func IsValidRegion(region string) bool {
	known := map[string]bool{
		"Africa": true, "Americas": true, "Asia": true,
		"Europe": true, "Oceania": true, "Antarctic": true,
	}
	return region == "" || known[region]
}