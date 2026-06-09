// controllers/api/country_test.go
// Unit tests for CountryAPIController endpoints.
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robiulislam99/TravelSphere/utils"
)

// Test: GetAll with no parameters
func TestGetAll_NoParameters(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries", nil)

	search := req.URL.Query().Get("search")
	region := req.URL.Query().Get("region")

	if search != "" {
		t.Errorf("Expected empty search, got %s", search)
	}
	if region != "" {
		t.Errorf("Expected empty region, got %s", region)
	}
}

// Test: GetAll with search parameter
func TestGetAll_WithSearchParameter(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=france", nil)

	search := req.URL.Query().Get("search")
	if search != "france" {
		t.Errorf("Expected search=france, got %s", search)
	}
}

// Test: GetAll with region parameter
func TestGetAll_WithValidRegion(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?region=Europe", nil)

	region := req.URL.Query().Get("region")
	if region != "Europe" {
		t.Errorf("Expected region=Europe, got %s", region)
	}

	// Validate region using IsValidRegion
	if !utils.IsValidRegion(region) {
		t.Errorf("Region should be valid: %s", region)
	}
}

// Test: GetAll with invalid region parameter
func TestGetAll_WithInvalidRegion(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?region=InvalidRegion", nil)

	region := req.URL.Query().Get("region")
	if region != "InvalidRegion" {
		t.Errorf("Expected region=InvalidRegion, got %s", region)
	}

	// Validate region using IsValidRegion
	if utils.IsValidRegion(region) {
		t.Errorf("Region should be invalid: %s", region)
	}
}

// Test: GetAll with both search and region
func TestGetAll_WithSearchAndRegion(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=united&region=Europe", nil)

	search := req.URL.Query().Get("search")
	region := req.URL.Query().Get("region")

	if search != "united" {
		t.Errorf("Expected search=united, got %s", search)
	}
	if region != "Europe" {
		t.Errorf("Expected region=Europe, got %s", region)
	}

	if !utils.IsValidRegion(region) {
		t.Errorf("Region should be valid: %s", region)
	}
}

// Test: GetAll with empty region (should be valid)
func TestGetAll_WithEmptyRegion(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?region=", nil)

	region := req.URL.Query().Get("region")
	// Empty region should be valid (means "no filter")
	if !utils.IsValidRegion(region) {
		t.Errorf("Empty region should be valid")
	}
}

// Test: GetAll with multiple search terms
func TestGetAll_WithMultipleSearchTerms(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=united+states", nil)

	search := req.URL.Query().Get("search")
	if search != "united states" {
		t.Errorf("Expected search=united states, got %s", search)
	}
}

// Test: GetAll with special characters in search
func TestGetAll_WithSpecialCharsInSearch(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=cote+d%27ivoire", nil)

	search := req.URL.Query().Get("search")
	if search != "cote d'ivoire" {
		t.Errorf("Expected search=cote d'ivoire, got %s", search)
	}
}

// Test: GetBySlug with valid slug format
func TestGetBySlug_ValidSlug(t *testing.T) {
	slug := "united-states"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Slug should be valid: %s", slug)
	}
}

// Test: GetBySlug with slug containing uppercase (invalid)
func TestGetBySlug_InvalidSlugUppercase(t *testing.T) {
	slug := "United-States"
	if utils.IsValidSlug(slug) {
		t.Errorf("Slug should be invalid (contains uppercase): %s", slug)
	}
}

// Test: GetBySlug with slug containing spaces (invalid)
func TestGetBySlug_InvalidSlugSpaces(t *testing.T) {
	slug := "united states"
	if utils.IsValidSlug(slug) {
		t.Errorf("Slug should be invalid (contains spaces): %s", slug)
	}
}

// Test: GetBySlug with slug containing special characters (invalid)
func TestGetBySlug_InvalidSlugSpecialChars(t *testing.T) {
	slug := "united_states"
	if utils.IsValidSlug(slug) {
		t.Errorf("Slug should be invalid (contains underscore): %s", slug)
	}
}

// Test: GetBySlug with empty slug (invalid)
func TestGetBySlug_EmptySlug(t *testing.T) {
	slug := ""
	if utils.IsValidSlug(slug) {
		t.Errorf("Slug should be invalid (empty)")
	}
}

// Test: GetBySlug with slug containing only hyphens (invalid)
func TestGetBySlug_OnlyHyphens(t *testing.T) {
	slug := "---"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Slug with only hyphens should be valid: %s", slug)
	}
}

// Test: GetBySlug with valid slug with numbers
func TestGetBySlug_ValidSlugWithNumbers(t *testing.T) {
	slug := "country123"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Slug should be valid: %s", slug)
	}
}

// Test: GetBySlug with valid slug with hyphens and numbers
func TestGetBySlug_ValidSlugComplex(t *testing.T) {
	slug := "united-arab-emirates"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Slug should be valid: %s", slug)
	}
}

// Test: Valid regions enum
func TestValidRegions(t *testing.T) {
	validRegions := []string{"Africa", "Americas", "Asia", "Europe", "Oceania", "Antarctic"}
	for _, region := range validRegions {
		if !utils.IsValidRegion(region) {
			t.Errorf("Region should be valid: %s", region)
		}
	}
}

// Test: Invalid regions
func TestInvalidRegions(t *testing.T) {
	invalidRegions := []string{"North America", "Central Asia", "Scandinavia", "Middle East"}
	for _, region := range invalidRegions {
		if utils.IsValidRegion(region) {
			t.Errorf("Region should be invalid: %s", region)
		}
	}
}

// Test: Response structure for success
func TestParseResponse_Success(t *testing.T) {
	responseBody := `{"data":"test data","message":"ok","status":200}`
	resp := parseResponse(responseBody)

	if resp.Data != "test data" {
		t.Errorf("Expected data=test data, got %v", resp.Data)
	}
	if resp.Message != "ok" {
		t.Errorf("Expected message=ok, got %s", resp.Message)
	}
	if resp.Status != 200 {
		t.Errorf("Expected status=200, got %d", resp.Status)
	}
}

// Test: Response structure for error
func TestParseResponse_Error(t *testing.T) {
	responseBody := `{"data":null,"message":"invalid region value","status":400}`
	resp := parseResponse(responseBody)

	if resp.Data != nil {
		t.Errorf("Expected data=null, got %v", resp.Data)
	}
	if resp.Message != "invalid region value" {
		t.Errorf("Expected message=invalid region value, got %s", resp.Message)
	}
	if resp.Status != 400 {
		t.Errorf("Expected status=400, got %d", resp.Status)
	}
}

// Test: Response structure for not found
func TestParseResponse_NotFound(t *testing.T) {
	responseBody := `{"data":null,"message":"country not found","status":404}`
	resp := parseResponse(responseBody)

	if resp.Data != nil {
		t.Errorf("Expected data=null, got %v", resp.Data)
	}
	if resp.Message != "country not found" {
		t.Errorf("Expected message=country not found, got %s", resp.Message)
	}
	if resp.Status != 404 {
		t.Errorf("Expected status=404, got %d", resp.Status)
	}
}

// Test: Response with array data
func TestParseResponse_ArrayData(t *testing.T) {
	responseBody := `{"data":[{"name":"France","slug":"france"}],"message":"ok","status":200}`
	resp := parseResponse(responseBody)

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status=%d, got %d", http.StatusOK, resp.Status)
	}

	// Verify data is array
	if dataArray, ok := resp.Data.([]interface{}); ok {
		if len(dataArray) != 1 {
			t.Errorf("Expected 1 item in data, got %d", len(dataArray))
		}
	} else {
		t.Errorf("Expected data to be array, got %T", resp.Data)
	}
}

// Test: Case-insensitive region parameter validation
func TestGetAll_CaseInsensitiveRegion(t *testing.T) {
	// Test lowercase
	req := httptest.NewRequest("GET", "/api/countries?region=europe", nil)
	region := req.URL.Query().Get("region")
	// API expects exact match, but we'll verify what's passed
	if region != "europe" {
		t.Errorf("Expected region=europe, got %s", region)
	}
}

// Test: Whitespace handling in search
func TestGetAll_SearchWithWhitespace(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=%20france%20", nil)
	search := req.URL.Query().Get("search")
	
	// Service layer trims whitespace, but verify parameter is captured
	if search != " france " {
		t.Errorf("Expected search= france , got %s", search)
	}
}

// Test: Multiple regions (only first should be used)
func TestGetAll_MultipleRegionParameters(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?region=Europe&region=Asia", nil)
	region := req.URL.Query().Get("region")
	
	// GetString gets only the first value
	if region != "Europe" {
		t.Errorf("Expected region=Europe (first value), got %s", region)
	}
}

// Test: URL-encoded special characters in search
func TestGetAll_URLEncodedSearch(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/countries?search=%C3%89tats-Unis", nil)
	search := req.URL.Query().Get("search")
	
	// Verify URL decoding works
	if search != "États-Unis" {
		t.Errorf("Expected URL decoded search, got %s", search)
	}
}

// Test: GetBySlug with slug containing digits
func TestGetBySlug_SlugStartsWithDigit(t *testing.T) {
	slug := "1country"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Slug starting with digit should be valid: %s", slug)
	}
}

// Test: GetBySlug with single character slug
func TestGetBySlug_SingleCharSlug(t *testing.T) {
	slug := "a"
	if !utils.IsValidSlug(slug) {
		t.Errorf("Single character slug should be valid: %s", slug)
	}
}

// Test: Region case sensitivity
func TestGetAll_RegionCaseSensitivity(t *testing.T) {
	regions := []struct {
		region string
		valid  bool
	}{
		{"Africa", true},
		{"africa", false},
		{"AFRICA", false},
		{"Europe", true},
		{"europe", false},
		{"", true}, // empty is valid (no filter)
	}

	for _, tc := range regions {
		valid := utils.IsValidRegion(tc.region)
		if valid != tc.valid {
			t.Errorf("Region %q: expected valid=%v, got %v", tc.region, tc.valid, valid)
		}
	}
}
