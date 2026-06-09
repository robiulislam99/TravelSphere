package services

import (
	"fmt"
	"strings"
	"testing"

	"github.com/robiulislam99/TravelSphere/models"
)

// TestNewWishlistService tests service initialization
func TestNewWishlistService(t *testing.T) {
	service := NewWishlistService()

	if service == nil {
		t.Error("WishlistService should not be nil")
	}
	if service.entries == nil {
		t.Error("WishlistService.entries should be initialized")
	}
}

// TestGetAll_EmptyWishlist tests retrieving entries from empty wishlist
func TestGetAll_EmptyWishlist(t *testing.T) {
	service := NewWishlistService()

	entries := service.GetAll("testuser")

	if entries == nil {
		t.Error("GetAll should return an empty slice, not nil")
	}
	if len(entries) != 0 {
		t.Errorf("GetAll on empty wishlist returned %d entries, want 0", len(entries))
	}
}

// TestGetAll_MultipleEntries tests retrieving multiple entries
func TestGetAll_MultipleEntries(t *testing.T) {
	service := NewWishlistService()

	// Create some entries
	req1 := &models.CreateWishlistRequest{
		CountryName: "France",
		Note:        "Visit Paris",
		Status:      string(models.StatusPlanned),
	}
	service.Create("testuser", req1)

	req2 := &models.CreateWishlistRequest{
		CountryName: "Italy",
		Note:        "Visit Rome",
		Status:      string(models.StatusPlanned),
	}
	service.Create("testuser", req2)

	entries := service.GetAll("testuser")

	if len(entries) != 2 {
		t.Errorf("GetAll returned %d entries, want 2", len(entries))
	}
}

// TestCreate_ValidEntry tests creating a valid entry
func TestCreate_ValidEntry(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Japan",
		Note:        "Beautiful country",
		Status:      string(models.StatusPlanned),
	}

	entry, err := service.Create("testuser", req)

	if err != nil {
		t.Errorf("Create returned error: %v", err)
	}
	if entry == nil {
		t.Error("Create should return a WishlistEntry")
	}
	if entry.CountryName != "Japan" {
		t.Errorf("Create().CountryName = %s, want Japan", entry.CountryName)
	}
	if entry.Note != "Beautiful country" {
		t.Errorf("Create().Note = %s, want Beautiful country", entry.Note)
	}
	if entry.Status != models.StatusPlanned {
		t.Errorf("Create().Status = %s, want %s", entry.Status, models.StatusPlanned)
	}
	if entry.ID == "" {
		t.Error("Create().ID should not be empty")
	}
	if entry.CreatedAt.IsZero() {
		t.Error("Create().CreatedAt should be set")
	}
}

// TestCreate_MissingCountryName tests validation of required field
func TestCreate_MissingCountryName(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "",
		Note:        "Some note",
		Status:      string(models.StatusPlanned),
	}

	entry, err := service.Create("testuser", req)

	if err == nil {
		t.Error("Create with empty country_name should return error")
	}
	if entry != nil {
		t.Error("Create with error should return nil entry")
	}
}

// TestCreate_InvalidStatus tests validation of invalid status
func TestCreate_InvalidStatus(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "France",
		Status:      "InvalidStatus",
	}

	entry, err := service.Create("testuser", req)

	if err == nil {
		t.Error("Create with invalid status should return error")
	}
	if entry != nil {
		t.Error("Create with error should return nil entry")
	}
}

// TestCreate_DefaultStatus tests that empty status defaults to Planned
func TestCreate_DefaultStatus(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Status:      "",
	}

	entry, err := service.Create("testuser", req)

	if err != nil {
		t.Errorf("Create should not error on empty status: %v", err)
	}
	if entry.Status != models.StatusPlanned {
		t.Errorf("Create with empty status should default to Planned, got %s", entry.Status)
	}
}

// TestGetByID_Found tests retrieving an entry by ID
func TestGetByID_Found(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Greece",
		Note:        "Visit islands",
		Status:      string(models.StatusPlanned),
	}
	created, _ := service.Create("testuser", req)

	retrieved, err := service.GetByID("testuser", created.ID)

	if err != nil {
		t.Errorf("GetByID returned error: %v", err)
	}
	if retrieved == nil {
		t.Error("GetByID should return an entry")
	}
	if retrieved.CountryName != "Greece" {
		t.Errorf("GetByID().CountryName = %s, want Greece", retrieved.CountryName)
	}
}

// TestGetByID_NotFound tests retrieving a non-existent entry
func TestGetByID_NotFound(t *testing.T) {
	service := NewWishlistService()

	retrieved, err := service.GetByID("testuser", "non-existent-id")

	if err != nil {
		t.Errorf("GetByID should not return error for missing entry: %v", err)
	}
	if retrieved != nil {
		t.Error("GetByID should return nil for non-existent entry")
	}
}

// TestUpdate_ValidEntry tests updating an entry
func TestUpdate_ValidEntry(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Portugal",
		Note:        "Old note",
		Status:      string(models.StatusPlanned),
	}
	created, _ := service.Create("testuser", req)

	updateReq := &models.UpdateWishlistRequest{
		Note:   "Updated note",
		Status: string(models.StatusVisited),
	}
	updated, err := service.Update("testuser", created.ID, updateReq)

	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}
	if updated.Note != "Updated note" {
		t.Errorf("Update().Note = %s, want Updated note", updated.Note)
	}
	if updated.Status != models.StatusVisited {
		t.Errorf("Update().Status = %s, want %s", updated.Status, models.StatusVisited)
	}
}

// TestUpdate_NotFound tests updating a non-existent entry
func TestUpdate_NotFound(t *testing.T) {
	service := NewWishlistService()

	updateReq := &models.UpdateWishlistRequest{
		Note:   "Some note",
		Status: string(models.StatusPlanned),
	}
	updated, err := service.Update("testuser", "non-existent-id", updateReq)

	if err == nil {
		t.Error("Update of non-existent entry should return error")
	}
	if updated != nil {
		t.Error("Update with error should return nil")
	}
}

// TestDelete_Success tests deleting an entry
func TestDelete_Success(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Turkey",
		Status:      string(models.StatusPlanned),
	}
	created, _ := service.Create("testuser", req)

	err := service.Delete("testuser", created.ID)

	if err != nil {
		t.Errorf("Delete returned error: %v", err)
	}

	// Verify entry is deleted
	retrieved, _ := service.GetByID("testuser", created.ID)
	if retrieved != nil {
		t.Error("Deleted entry should not be retrievable")
	}
}

// TestDelete_NotFound tests deleting a non-existent entry
func TestDelete_NotFound(t *testing.T) {
	service := NewWishlistService()

	err := service.Delete("testuser", "non-existent-id")

	if err == nil {
		t.Error("Delete of non-existent entry should return error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Delete error message should mention 'not found', got: %v", err)
	}
}

// TestIsDuplicate_True tests duplicate detection
func TestIsDuplicate_True(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Germany",
		Status:      string(models.StatusPlanned),
	}
	service.Create("testuser", req)

	isDuplicate := service.IsDuplicate("testuser", "Germany")

	if !isDuplicate {
		t.Error("IsDuplicate should return true for existing country")
	}
}

// TestIsDuplicate_False tests when not a duplicate
func TestIsDuplicate_False(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Germany",
		Status:      string(models.StatusPlanned),
	}
	service.Create("testuser", req)

	isDuplicate := service.IsDuplicate("testuser", "France")

	if isDuplicate {
		t.Error("IsDuplicate should return false for non-existing country")
	}
}

// TestIsDuplicate_CaseInsensitive tests case-insensitive duplicate detection
func TestIsDuplicate_CaseInsensitive(t *testing.T) {
	service := NewWishlistService()

	req := &models.CreateWishlistRequest{
		CountryName: "Germany",
		Status:      string(models.StatusPlanned),
	}
	service.Create("testuser", req)

	tests := []struct {
		countryName string
		expected    bool
	}{
		{"germany", true},
		{"GERMANY", true},
		{"Germany", true},
		{"France", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("IsDuplicate_%s", tt.countryName), func(t *testing.T) {
			result := service.IsDuplicate("testuser", tt.countryName)
			if result != tt.expected {
				t.Errorf("IsDuplicate(%q) = %v, want %v", tt.countryName, result, tt.expected)
			}
		})
	}
}

// TestUserIsolation tests that wishlist entries are isolated per user
func TestUserIsolation(t *testing.T) {
	service := NewWishlistService()

	// Add entry for user1
	req1 := &models.CreateWishlistRequest{
		CountryName: "Canada",
		Status:      string(models.StatusPlanned),
	}
	service.Create("user1", req1)

	// Add entry for user2
	req2 := &models.CreateWishlistRequest{
		CountryName: "Mexico",
		Status:      string(models.StatusPlanned),
	}
	service.Create("user2", req2)

	// Get entries for each user
	user1Entries := service.GetAll("user1")
	user2Entries := service.GetAll("user2")

	if len(user1Entries) != 1 {
		t.Errorf("user1 should have 1 entry, got %d", len(user1Entries))
	}
	if len(user2Entries) != 1 {
		t.Errorf("user2 should have 1 entry, got %d", len(user2Entries))
	}

	if user1Entries[0].CountryName != "Canada" {
		t.Errorf("user1 entry should be Canada, got %s", user1Entries[0].CountryName)
	}
	if user2Entries[0].CountryName != "Mexico" {
		t.Errorf("user2 entry should be Mexico, got %s", user2Entries[0].CountryName)
	}
}

// TestGetAll_SortingNewestFirst tests that entries are sorted newest first
func TestGetAll_SortingNewestFirst(t *testing.T) {
	service := NewWishlistService()

	// Create first entry
	req1 := &models.CreateWishlistRequest{
		CountryName: "Argentina",
		Status:      string(models.StatusPlanned),
	}
	entry1, _ := service.Create("testuser", req1)

	// Create second entry
	req2 := &models.CreateWishlistRequest{
		CountryName: "Brazil",
		Status:      string(models.StatusPlanned),
	}
	entry2, _ := service.Create("testuser", req2)

	entries := service.GetAll("testuser")

	// Newest (entry2) should come first
	if entries[0].ID != entry2.ID {
		t.Error("Most recent entry should be first in GetAll results")
	}
	if entries[1].ID != entry1.ID {
		t.Error("Oldest entry should be last in GetAll results")
	}
}
