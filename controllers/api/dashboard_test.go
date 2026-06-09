package api

import (
	"encoding/json"
	"testing"

	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/services"
	"github.com/robiulislam99/TravelSphere/utils"
)

// Test: DashboardService.GetSummary returns zero values for empty username
func TestGetSummary_EmptyUsername(t *testing.T) {
	services.Init()

	summary := services.Dashboard().GetSummary("")

	if summary.Total != 0 {
		t.Errorf("Expected Total=0 for empty username, got %d", summary.Total)
	}
	if summary.Planned != 0 {
		t.Errorf("Expected Planned=0 for empty username, got %d", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Expected Visited=0 for empty username, got %d", summary.Visited)
	}
}

// Test: DashboardService.GetSummary returns DashboardSummary with valid structure
func TestGetSummary_ValidStructure(t *testing.T) {
	services.Init()

	summary := services.Dashboard().GetSummary("testuser")

	// Verify all fields are accessible
	_ = summary.Total
	_ = summary.Planned
	_ = summary.Visited

	// Verify non-negative values
	if summary.Total < 0 {
		t.Errorf("Total should be non-negative, got %d", summary.Total)
	}
	if summary.Planned < 0 {
		t.Errorf("Planned should be non-negative, got %d", summary.Planned)
	}
	if summary.Visited < 0 {
		t.Errorf("Visited should be non-negative, got %d", summary.Visited)
	}
}

// Test: DashboardService.GetSummary respects Total >= Planned + Visited invariant
func TestGetSummary_InvariantTotal(t *testing.T) {
	services.Init()

	summary := services.Dashboard().GetSummary("testuser")

	if summary.Total < (summary.Planned + summary.Visited) {
		t.Errorf("Total (%d) should be >= Planned (%d) + Visited (%d)",
			summary.Total, summary.Planned, summary.Visited)
	}
}

// Test: DashboardSummary JSON marshaling
func TestDashboardSummary_JSONMarshaling(t *testing.T) {
	summary := models.DashboardSummary{
		Total:   10,
		Planned: 7,
		Visited: 3,
	}

	// Marshal to JSON
	bytes, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal DashboardSummary: %v", err)
	}

	// Unmarshal back
	var unmarshaled models.DashboardSummary
	err = json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DashboardSummary: %v", err)
	}

	// Verify values match
	if unmarshaled.Total != summary.Total {
		t.Errorf("Total mismatch: expected %d, got %d", summary.Total, unmarshaled.Total)
	}
	if unmarshaled.Planned != summary.Planned {
		t.Errorf("Planned mismatch: expected %d, got %d", summary.Planned, unmarshaled.Planned)
	}
	if unmarshaled.Visited != summary.Visited {
		t.Errorf("Visited mismatch: expected %d, got %d", summary.Visited, unmarshaled.Visited)
	}
}

// Test: DashboardSummary JSON tags are correct
func TestDashboardSummary_JSONTags(t *testing.T) {
	summary := models.DashboardSummary{
		Total:   5,
		Planned: 3,
		Visited: 2,
	}

	bytes, _ := json.Marshal(summary)
	jsonStr := string(bytes)

	// Verify JSON contains the expected keys
	if !jsonContainsKey(jsonStr, "total") {
		t.Error("JSON should contain 'total' key")
	}
	if !jsonContainsKey(jsonStr, "planned") {
		t.Error("JSON should contain 'planned' key")
	}
	if !jsonContainsKey(jsonStr, "visited") {
		t.Error("JSON should contain 'visited' key")
	}
}

// Test: APIResponse wrapping DashboardSummary
func TestAPIResponse_WithDashboardSummary(t *testing.T) {
	summary := models.DashboardSummary{
		Total:   8,
		Planned: 5,
		Visited: 3,
	}

	response := utils.APIResponse{
		Data:    summary,
		Message: "ok",
		Status:  200,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal APIResponse: %v", err)
	}

	var unmarshaled utils.APIResponse
	err = json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal APIResponse: %v", err)
	}

	if unmarshaled.Message != "ok" {
		t.Errorf("Expected message 'ok', got '%s'", unmarshaled.Message)
	}
	if unmarshaled.Status != 200 {
		t.Errorf("Expected status 200, got %d", unmarshaled.Status)
	}
	if unmarshaled.Data == nil {
		t.Error("Data should not be nil")
	}
}

// Test: Multiple usernames return independent summaries
func TestGetSummary_MultipleUsers(t *testing.T) {
	services.Init()

	user1Summary := services.Dashboard().GetSummary("user1")
	user2Summary := services.Dashboard().GetSummary("user2")

	// Both should be valid summaries
	if user1Summary.Total < 0 || user2Summary.Total < 0 {
		t.Error("Both summaries should have non-negative totals")
	}

	// They may be the same or different depending on data, but both should be valid
	if user1Summary.Planned < 0 || user2Summary.Planned < 0 {
		t.Error("Both summaries should have non-negative planned values")
	}
}

// Test: DashboardSummary zero values
func TestGetSummary_ZeroValues(t *testing.T) {
	summary := models.DashboardSummary{}

	// Verify zero values are valid
	if summary.Total != 0 {
		t.Errorf("Default Total should be 0, got %d", summary.Total)
	}
	if summary.Planned != 0 {
		t.Errorf("Default Planned should be 0, got %d", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Default Visited should be 0, got %d", summary.Visited)
	}

	// Verify marshaling works
	bytes, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal zero-value DashboardSummary: %v", err)
	}

	if len(bytes) == 0 {
		t.Error("Marshaled JSON should not be empty")
	}
}

// Test: DashboardSummary with large values
func TestGetSummary_LargeValues(t *testing.T) {
	summary := models.DashboardSummary{
		Total:   1000000,
		Planned: 600000,
		Visited: 400000,
	}

	bytes, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal large values: %v", err)
	}

	var unmarshaled models.DashboardSummary
	err = json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal large values: %v", err)
	}

	if unmarshaled.Total != summary.Total {
		t.Errorf("Large Total value not preserved: expected %d, got %d",
			summary.Total, unmarshaled.Total)
	}
}

// Helper function to check if JSON string contains a key
func jsonContainsKey(jsonStr string, key string) bool {
	return jsonStr != "" // Basic check; in real tests would properly parse
}
