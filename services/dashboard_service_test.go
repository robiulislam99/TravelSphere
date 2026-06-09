package services

import (
	"testing"

	"github.com/robiulislam99/TravelSphere/models"
)

// TestNewDashboardService tests service initialization
func TestNewDashboardService(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	if dashboardService == nil {
		t.Error("DashboardService should not be nil")
	}
	if dashboardService.wishlist != wishlistService {
		t.Error("DashboardService.wishlist should reference the provided WishlistService")
	}
}

// TestGetSummary_EmptyWishlist tests summary with no entries
func TestGetSummary_EmptyWishlist(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	summary := dashboardService.GetSummary("testuser")

	if summary.Total != 0 {
		t.Errorf("Empty wishlist Total = %d, want 0", summary.Total)
	}
	if summary.Planned != 0 {
		t.Errorf("Empty wishlist Planned = %d, want 0", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Empty wishlist Visited = %d, want 0", summary.Visited)
	}
}

// TestGetSummary_WithPlannedEntry tests summary with planned entries
func TestGetSummary_WithPlannedEntries(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	// Create some planned entries
	req1 := &models.CreateWishlistRequest{
		CountryName: "France",
		Note:        "Visit in summer",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("testuser", req1)

	req2 := &models.CreateWishlistRequest{
		CountryName: "Italy",
		Note:        "Beautiful country",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("testuser", req2)

	summary := dashboardService.GetSummary("testuser")

	if summary.Total != 2 {
		t.Errorf("Summary Total = %d, want 2", summary.Total)
	}
	if summary.Planned != 2 {
		t.Errorf("Summary Planned = %d, want 2", summary.Planned)
	}
	if summary.Visited != 0 {
		t.Errorf("Summary Visited = %d, want 0", summary.Visited)
	}
}

// TestGetSummary_WithVisitedEntry tests summary with visited entries
func TestGetSummary_WithVisitedEntries(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	// Create a visited entry
	req := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Note:        "Already visited",
		Status:      string(models.StatusVisited),
	}
	wishlistService.Create("testuser", req)

	summary := dashboardService.GetSummary("testuser")

	if summary.Total != 1 {
		t.Errorf("Summary Total = %d, want 1", summary.Total)
	}
	if summary.Visited != 1 {
		t.Errorf("Summary Visited = %d, want 1", summary.Visited)
	}
	if summary.Planned != 0 {
		t.Errorf("Summary Planned = %d, want 0", summary.Planned)
	}
}

// TestGetSummary_MixedStatuses tests summary with both planned and visited
func TestGetSummary_MixedStatuses(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	// Create mixed entries
	planned := &models.CreateWishlistRequest{
		CountryName: "Japan",
		Note:        "Plan to visit",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("testuser", planned)

	visited := &models.CreateWishlistRequest{
		CountryName: "Thailand",
		Note:        "Already visited",
		Status:      string(models.StatusVisited),
	}
	wishlistService.Create("testuser", visited)

	planned2 := &models.CreateWishlistRequest{
		CountryName: "Greece",
		Note:        "Another planned trip",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("testuser", planned2)

	summary := dashboardService.GetSummary("testuser")

	if summary.Total != 3 {
		t.Errorf("Summary Total = %d, want 3", summary.Total)
	}
	if summary.Planned != 2 {
		t.Errorf("Summary Planned = %d, want 2", summary.Planned)
	}
	if summary.Visited != 1 {
		t.Errorf("Summary Visited = %d, want 1", summary.Visited)
	}
}

// TestGetSummary_UserIsolation tests that dashboard is scoped per user
func TestGetSummary_UserIsolation(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	// Add entries for user1
	req1 := &models.CreateWishlistRequest{
		CountryName: "France",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("user1", req1)

	req2 := &models.CreateWishlistRequest{
		CountryName: "Italy",
		Status:      string(models.StatusPlanned),
	}
	wishlistService.Create("user1", req2)

	// Add entries for user2
	req3 := &models.CreateWishlistRequest{
		CountryName: "Spain",
		Status:      string(models.StatusVisited),
	}
	wishlistService.Create("user2", req3)

	summary1 := dashboardService.GetSummary("user1")
	summary2 := dashboardService.GetSummary("user2")

	if summary1.Total != 2 {
		t.Errorf("user1 Summary Total = %d, want 2", summary1.Total)
	}
	if summary1.Planned != 2 {
		t.Errorf("user1 Summary Planned = %d, want 2", summary1.Planned)
	}

	if summary2.Total != 1 {
		t.Errorf("user2 Summary Total = %d, want 1", summary2.Total)
	}
	if summary2.Visited != 1 {
		t.Errorf("user2 Summary Visited = %d, want 1", summary2.Visited)
	}
}

// TestGetSummary_DefaultStatus tests entries with default status
func TestGetSummary_DefaultStatus(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	// Create entry without explicit status (should default to Planned)
	req := &models.CreateWishlistRequest{
		CountryName: "Germany",
		Note:        "Test default",
		Status:      "",
	}
	wishlistService.Create("testuser", req)

	summary := dashboardService.GetSummary("testuser")

	if summary.Total != 1 {
		t.Errorf("Summary Total = %d, want 1", summary.Total)
	}
	if summary.Planned != 1 {
		t.Errorf("Default status should be Planned, got Planned=%d, Visited=%d", summary.Planned, summary.Visited)
	}
}
