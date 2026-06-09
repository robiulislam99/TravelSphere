// controllers/api/dashboard.go
package api

import (
    "github.com/beego/beego/v2/server/web"
    "github.com/robiulislam99/TravelSphere/services"
    "github.com/robiulislam99/TravelSphere/utils"
)

type DashboardAPIController struct {
    web.Controller
}

func (c *DashboardAPIController) Summary() {
    username := ""
    if u := c.GetSession("username"); u != nil {
        username = u.(string)
    }

    summary := services.Dashboard().GetSummary(username)
    utils.SendSuccess(&c.Controller, summary)
}