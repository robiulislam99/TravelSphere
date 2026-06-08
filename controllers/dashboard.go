// controllers/dashboard.go
// DashboardController handles GET /dashboard (protected SSR page).
package controllers

import "github.com/robiulislam99/TravelSphere/services"

type DashboardController struct {
	BaseController
}

func (c *DashboardController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"]      = "Dashboard"
	c.Data["ActivePage"] = "dashboard"
}

func (c *DashboardController) Get() {
	summary := services.Dashboard().GetSummary()
	entries := services.Wishlist().GetAll()

	c.Data["Summary"]         = summary
	c.Data["WishlistEntries"] = entries

	c.Layout  = "layout/base.tpl"
	c.TplName = "pages/dashboard.tpl"
}