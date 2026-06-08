// controllers/api/dashboard.go
// DashboardAPIController handles GET /api/dashboard/summary.
// Full implementation in Phase 6.
package api

import "github.com/beego/beego/v2/server/web"

// DashboardAPIController serves the dashboard summary JSON.
type DashboardAPIController struct {
	web.Controller
}

// Summary handles GET /api/dashboard/summary
func (c *DashboardAPIController) Summary() {
	c.Data["json"] = map[string]interface{}{
		"total":   0,
		"planned": 0,
		"visited": 0,
	}
	c.ServeJSON()
}