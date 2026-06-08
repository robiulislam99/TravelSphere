// controllers/dashboard.go
// DashboardController handles GET /dashboard (protected route).
// Full implementation with DashboardService comes in Phase 5.
package controllers

// DashboardController handles the dashboard SSR page.
type DashboardController struct {
	BaseController
}

// Prepare sets page-level template data.
func (c *DashboardController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"] = "Dashboard"
	c.Data["ActivePage"] = "dashboard"
}

// Get renders the dashboard page.
// TODO (Phase 5): fetch summary and entries from DashboardService + WishlistService.
func (c *DashboardController) Get() {
	c.Data["Summary"] = map[string]int{
		"Total":   0,
		"Planned": 0,
		"Visited": 0,
	}
	c.Data["WishlistEntries"] = []interface{}{}

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/dashboard.tpl"
}