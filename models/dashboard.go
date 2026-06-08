// models/dashboard.go
// DashboardSummary holds aggregated wishlist statistics for the dashboard page
// and the /api/dashboard/summary JSON endpoint.
package models

// DashboardSummary contains counts used by the dashboard page and AJAX refresh.
type DashboardSummary struct {
	Total   int `json:"total"`
	Planned int `json:"planned"`
	Visited int `json:"visited"`
}