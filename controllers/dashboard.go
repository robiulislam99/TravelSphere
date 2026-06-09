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
    username := ""
    if u := c.GetSession("username"); u != nil {
        username = u.(string)
    }

    summary := services.Dashboard().GetSummary(username)
    entries := services.Wishlist().GetAll(username)

    c.Data["Summary"]         = summary
    c.Data["WishlistEntries"] = entries
    c.Data["PageScript"]      = "dashboard.js"
    c.Layout                  = "layout/base.tpl"
    c.TplName                 = "pages/dashboard.tpl"
}