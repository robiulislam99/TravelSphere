// services/dashboard_service.go
// DashboardService computes summary statistics from the WishlistService.
// Keeps dashboard logic out of the controller.
package services

import "github.com/robiulislam99/TravelSphere/models"

// DashboardService aggregates wishlist data for the dashboard page.
type DashboardService struct {
	wishlist *WishlistService
}

// NewDashboardService creates a DashboardService backed by the given WishlistService.
func NewDashboardService(ws *WishlistService) *DashboardService {
	return &DashboardService{wishlist: ws}
}

// GetSummary returns total, planned, and visited counts for the current wishlist.
func (s *DashboardService) GetSummary() models.DashboardSummary {
	entries := s.wishlist.GetAll()

	summary := models.DashboardSummary{Total: len(entries)}
	for _, e := range entries {
		switch e.Status {
		case models.StatusPlanned:
			summary.Planned++
		case models.StatusVisited:
			summary.Visited++
		}
	}
	return summary
}