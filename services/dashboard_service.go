package services

import "github.com/robiulislam99/TravelSphere/models"

type DashboardService struct {
    wishlist *WishlistService
}

func NewDashboardService(ws *WishlistService) *DashboardService {
    return &DashboardService{wishlist: ws}
}

// GetSummary now takes username so it only counts that user's entries.
func (s *DashboardService) GetSummary(username string) models.DashboardSummary {
    entries := s.wishlist.GetAll(username)

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