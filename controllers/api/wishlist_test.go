package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/robiulislam99/TravelSphere/models"
	"github.com/robiulislam99/TravelSphere/services"
)

// ========== HTTP Request/Response Tests ==========

// Test: HTTP GET request with no parameters
func TestWishlistHTTP_GetAllNoParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/wishlist", nil)

	if req.Method != "GET" {
		t.Errorf("Expected GET method, got %s", req.Method)
	}
	if req.URL.Path != "/api/wishlist" {
		t.Errorf("Expected path /api/wishlist, got %s", req.URL.Path)
	}
}

// Test: HTTP POST with valid JSON
func TestWishlistHTTP_PostValidJSON(t *testing.T) {
	reqBody := models.CreateWishlistRequest{
		CountryName: "Germany",
		Note:        "Test",
		Status:      string(models.StatusPlanned),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/wishlist", bytes.NewReader(body))

	if req.Method != "POST" {
		t.Errorf("Expected POST method, got %s", req.Method)
	}

	var parsed models.CreateWishlistRequest
	json.NewDecoder(req.Body).Decode(&parsed)

	if parsed.CountryName != "Germany" {
		t.Errorf("Expected country Germany, got %s", parsed.CountryName)
	}
}

// Test: HTTP PUT with ID parameter
func TestWishlistHTTP_PutWithIDParam(t *testing.T) {
	reqBody := models.UpdateWishlistRequest{
		Note:   "Updated",
		Status: string(models.StatusVisited),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/wishlist/test-id-123", bytes.NewReader(body))

	if req.Method != "PUT" {
		t.Errorf("Expected PUT method, got %s", req.Method)
	}

	// Extract ID from path
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) > 0 {
		extractedID := pathParts[len(pathParts)-1]
		if extractedID != "test-id-123" {
			t.Errorf("Expected ID test-id-123, got %s", extractedID)
		}
	}
}

// Test: HTTP DELETE with ID parameter
func TestWishlistHTTP_DeleteWithIDParam(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/wishlist/test-id-456", nil)

	if req.Method != "DELETE" {
		t.Errorf("Expected DELETE method, got %s", req.Method)
	}

	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) > 0 {
		extractedID := pathParts[len(pathParts)-1]
		if extractedID != "test-id-456" {
			t.Errorf("Expected ID test-id-456, got %s", extractedID)
		}
	}
}

// ========== Service Layer Tests ==========

// Test: GetAll returns empty list initially
func TestWishlistService_GetAll_Empty(t *testing.T) {
	svc := services.NewWishlistService()

	entries := svc.GetAll("testuser")

	if entries == nil {
		t.Errorf("Expected empty slice, got nil")
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}

// Test: GetAll returns created entries
func TestWishlistService_GetAll_WithEntries(t *testing.T) {
	svc := services.NewWishlistService()

	// Create some entries
	req1 := &models.CreateWishlistRequest{
		CountryName: "France",
		Note:        "Visit Paris",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("testuser", req1)

	req2 := &models.CreateWishlistRequest{
		CountryName: "Italy",
		Note:        "Visit Rome",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("testuser", req2)

	entries := svc.GetAll("testuser")

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

// Test: GetAll isolates entries per user
func TestWishlistService_GetAll_UserIsolation(t *testing.T) {
	svc := services.NewWishlistService()

	// User 1 creates entries
	req1 := &models.CreateWishlistRequest{
		CountryName: "Japan",
		Note:        "User 1",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("user1", req1)

	// User 2 creates entry
	req2 := &models.CreateWishlistRequest{
		CountryName: "Brazil",
		Note:        "User 2",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("user2", req2)

	entries1 := svc.GetAll("user1")
	entries2 := svc.GetAll("user2")

	if len(entries1) != 1 {
		t.Errorf("Expected user1 to have 1 entry, got %d", len(entries1))
	}
	if len(entries2) != 1 {
		t.Errorf("Expected user2 to have 1 entry, got %d", len(entries2))
	}

	if entries1[0].CountryName != "Japan" {
		t.Errorf("Expected user1 entry to be Japan, got %s", entries1[0].CountryName)
	}
	if entries2[0].CountryName != "Brazil" {
		t.Errorf("Expected user2 entry to be Brazil, got %s", entries2[0].CountryName)
	}
}

// Test: Create with valid request
func TestWishlistService_Create_Valid(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Japan",
		Note:        "Experience Tokyo",
		Status:      string(models.StatusPlanned),
	}

	entry, err := svc.Create("testuser", req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if entry == nil {
		t.Fatalf("Expected non-nil entry")
	}
	if entry.ID == "" {
		t.Errorf("Expected entry to have an ID")
	}
	if entry.CountryName != "Japan" {
		t.Errorf("Expected country name Japan, got %s", entry.CountryName)
	}
	if entry.Status != models.StatusPlanned {
		t.Errorf("Expected status Planned, got %s", entry.Status)
	}
}

// Test: Create with missing country_name
func TestWishlistService_Create_MissingCountryName(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "",
		Note:        "Some note",
		Status:      string(models.StatusPlanned),
	}

	entry, err := svc.Create("testuser", req)

	if err == nil {
		t.Errorf("Expected error for missing country_name")
	}
	if entry != nil {
		t.Errorf("Expected nil entry for invalid request")
	}
	if err.Error() != "country_name is required" {
		t.Errorf("Expected 'country_name is required' error, got %s", err.Error())
	}
}

// Test: Create with invalid status
func TestWishlistService_Create_InvalidStatus(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Italy",
		Note:        "Visit Rome",
		Status:      "InvalidStatus",
	}

	entry, err := svc.Create("testuser", req)

	if err == nil {
		t.Errorf("Expected error for invalid status")
	}
	if entry != nil {
		t.Errorf("Expected nil entry for invalid request")
	}
	if err.Error() != "status must be 'Planned' or 'Visited'" {
		t.Errorf("Expected status validation error, got %s", err.Error())
	}
}

// Test: Create with empty status defaults to Planned
func TestWishlistService_Create_EmptyStatusDefaultsToPlanned(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Germany",
		Note:        "Visit Berlin",
		Status:      "",
	}

	entry, err := svc.Create("testuser", req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if entry == nil {
		t.Fatalf("Expected non-nil entry")
	}
	if entry.Status != models.StatusPlanned {
		t.Errorf("Expected default status Planned, got %s", entry.Status)
	}
}

// Test: Create with empty note is allowed
func TestWishlistService_Create_EmptyNote(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Mexico",
		Note:        "",
		Status:      string(models.StatusPlanned),
	}

	entry, err := svc.Create("testuser", req)

	if err != nil {
		t.Errorf("Expected no error for empty note, got %v", err)
	}
	if entry == nil {
		t.Fatalf("Expected non-nil entry")
	}
}

// Test: IsDuplicate detects duplicate countries for same user
func TestWishlistService_IsDuplicate_DetectsDuplicate(t *testing.T) {
	svc := services.NewWishlistService()

	// Create first entry
	req1 := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Note:        "First",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("sameuser", req1)

	// Check duplicate before creating second
	isDup := svc.IsDuplicate("sameuser", "Spain")

	if !isDup {
		t.Errorf("Expected IsDuplicate to return true for existing country")
	}

	// Create second entry with same country
	req2 := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Note:        "Second",
		Status:      string(models.StatusVisited),
	}
	entry, err := svc.Create("sameuser", req2)

	if err != nil {
		t.Errorf("Expected create to succeed even with duplicate, got %v", err)
	}
	if entry == nil {
		t.Fatalf("Expected non-nil entry")
	}
}

// Test: IsDuplicate is case-insensitive
func TestWishlistService_IsDuplicate_CaseInsensitive(t *testing.T) {
	svc := services.NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "FRANCE",
		Note:        "Test",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("testuser", req)

	// Check with different case
	isDup := svc.IsDuplicate("testuser", "france")

	if !isDup {
		t.Errorf("Expected IsDuplicate to be case-insensitive")
	}
}

// Test: IsDuplicate not triggered across users
func TestWishlistService_IsDuplicate_UserIsolation(t *testing.T) {
	svc := services.NewWishlistService()

	// User 1 creates entry
	req1 := &models.CreateWishlistRequest{
		CountryName: "Canada",
		Note:        "User 1",
		Status:      string(models.StatusPlanned),
	}
	svc.Create("user1", req1)

	// User 2 checks - should not see User 1's entry as duplicate
	isDup := svc.IsDuplicate("user2", "Canada")

	if isDup {
		t.Errorf("Expected IsDuplicate to be false for different user")
	}
}

// Test: Update with valid request
func TestWishlistService_Update_Valid(t *testing.T) {
	svc := services.NewWishlistService()

	// Create an entry first
	createReq := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Note:        "Original note",
		Status:      string(models.StatusPlanned),
	}
	entry, _ := svc.Create("testuser", createReq)
	entryID := entry.ID

	// Update it
	updateReq := &models.UpdateWishlistRequest{
		Note:   "Updated note",
		Status: string(models.StatusVisited),
	}
	updatedEntry, err := svc.Update("testuser", entryID, updateReq)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updatedEntry == nil {
		t.Fatalf("Expected non-nil entry")
	}
	if updatedEntry.Note != "Updated note" {
		t.Errorf("Expected note 'Updated note', got '%s'", updatedEntry.Note)
	}
	if updatedEntry.Status != models.StatusVisited {
		t.Errorf("Expected status Visited, got %s", updatedEntry.Status)
	}
}

// Test: Update non-existent entry
func TestWishlistService_Update_NotFound(t *testing.T) {
	svc := services.NewWishlistService()

	updateReq := &models.UpdateWishlistRequest{
		Note:   "Some note",
		Status: string(models.StatusVisited),
	}

	entry, err := svc.Update("testuser", "nonexistent", updateReq)

	if err == nil {
		t.Errorf("Expected error for non-existent entry")
	}
	if entry != nil {
		t.Errorf("Expected nil entry for non-existent ID")
	}
	if err.Error() != "wishlist entry not found" {
		t.Errorf("Expected 'wishlist entry not found' error, got %s", err.Error())
	}
}

// Test: Update with invalid status
func TestWishlistService_Update_InvalidStatus(t *testing.T) {
	svc := services.NewWishlistService()

	// Create an entry first
	createReq := &models.CreateWishlistRequest{
		CountryName: "Canada",
		Note:        "Original",
		Status:      string(models.StatusPlanned),
	}
	entry, _ := svc.Create("testuser", createReq)

	// Try to update with invalid status
	updateReq := &models.UpdateWishlistRequest{
		Note:   "Updated",
		Status: "BadStatus",
	}

	updatedEntry, err := svc.Update("testuser", entry.ID, updateReq)

	if err == nil {
		t.Errorf("Expected error for invalid status")
	}
	if updatedEntry != nil {
		t.Errorf("Expected nil entry for invalid status")
	}
}

// Test: Update with empty status is rejected
func TestWishlistService_Update_EmptyStatus(t *testing.T) {
	svc := services.NewWishlistService()

	// Create an entry first
	createReq := &models.CreateWishlistRequest{
		CountryName: "Poland",
		Note:        "Original",
		Status:      string(models.StatusPlanned),
	}
	entry, _ := svc.Create("testuser", createReq)

	// Try to update with empty status
	updateReq := &models.UpdateWishlistRequest{
		Note:   "Updated",
		Status: "",
	}

	updatedEntry, err := svc.Update("testuser", entry.ID, updateReq)

	if err == nil {
		t.Errorf("Expected error for empty status")
	}
	if updatedEntry != nil {
		t.Errorf("Expected nil entry for empty status")
	}
}

// Test: Update only changes the specified fields
func TestWishlistService_Update_PartialUpdate(t *testing.T) {
	svc := services.NewWishlistService()

	// Create an entry first
	createReq := &models.CreateWishlistRequest{
		CountryName: "Greece",
		Note:        "Original note",
		Status:      string(models.StatusPlanned),
	}
	entry, _ := svc.Create("testuser", createReq)
	originalCountry := entry.CountryName

	// Update only note and status (country should remain the same)
	updateReq := &models.UpdateWishlistRequest{
		Note:   "New note",
		Status: string(models.StatusVisited),
	}
	updatedEntry, _ := svc.Update("testuser", entry.ID, updateReq)

	if updatedEntry.CountryName != originalCountry {
		t.Errorf("Expected country to remain %s, got %s", originalCountry, updatedEntry.CountryName)
	}
}

// Test: Delete existing entry
func TestWishlistService_Delete_Success(t *testing.T) {
	svc := services.NewWishlistService()

	// Create an entry first
	createReq := &models.CreateWishlistRequest{
		CountryName: "Australia",
		Note:        "To delete",
		Status:      string(models.StatusPlanned),
	}
	entry, _ := svc.Create("testuser", createReq)
	entryID := entry.ID

	// Delete it
	err := svc.Delete("testuser", entryID)

	if err != nil {
		t.Errorf("Expected no error on delete, got %v", err)
	}

	// Verify it's gone
	retrieved, _ := svc.GetByID("testuser", entryID)
	if retrieved != nil {
		t.Errorf("Expected entry to be deleted, but still found it")
	}
}

// Test: Delete non-existent entry
func TestWishlistService_Delete_NotFound(t *testing.T) {
	svc := services.NewWishlistService()

	err := svc.Delete("testuser", "nonexistent")

	if err == nil {
		t.Errorf("Expected error for non-existent entry")
	}
	if err.Error() != "wishlist entry not found" {
		t.Errorf("Expected 'wishlist entry not found' error, got %s", err.Error())
	}
}

// Test: Delete only affects specified user's entry
func TestWishlistService_Delete_UserIsolation(t *testing.T) {
	svc := services.NewWishlistService()

	// User 1 creates entry
	req1 := &models.CreateWishlistRequest{
		CountryName: "Thailand",
		Note:        "User 1",
		Status:      string(models.StatusPlanned),
	}
	entry1, _ := svc.Create("user1", req1)

	// User 2 creates entry
	req2 := &models.CreateWishlistRequest{
		CountryName: "Vietnam",
		Note:        "User 2",
		Status:      string(models.StatusPlanned),
	}
	entry2, _ := svc.Create("user2", req2)

	// Delete User 1's entry
	svc.Delete("user1", entry1.ID)

	// User 1 should have 0 entries
	entries1 := svc.GetAll("user1")
	if len(entries1) != 0 {
		t.Errorf("Expected user1 to have 0 entries after delete, got %d", len(entries1))
	}

	// User 2 should still have 1 entry
	entries2 := svc.GetAll("user2")
	if len(entries2) != 1 {
		t.Errorf("Expected user2 to still have 1 entry, got %d", len(entries2))
	}
	if entries2[0].ID != entry2.ID {
		t.Errorf("Expected user2's entry to be preserved")
	}
}
