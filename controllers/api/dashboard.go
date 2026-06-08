// controllers/api/dashboard.go
// DashboardAPIController — GET /api/dashboard/summary
package api

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/services"
	"github.com/robiulislam99/TravelSphere/utils"
)

type DashboardAPIController struct {
	web.Controller
}

// Summary handles GET /api/dashboard/summary
func (c *DashboardAPIController) Summary() {
	summary := services.Dashboard().GetSummary()
	utils.SendSuccess(&c.Controller, summary)
}