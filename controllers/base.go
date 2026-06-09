// controllers/base.go
// BaseController embedded by all SSR page controllers.
// Prepare() sets common template data available to every page.
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/services"
)

type BaseController struct {
	web.Controller
}

// Prepare runs before every action. Sets defaults every template can rely on.
func (c *BaseController) Prepare() {
	// Ensure services are ready (idempotent)
	services.Init()

	c.Data["Title"]           = "TravelSphere"
	c.Data["MetaDescription"] = "Discover destinations, explore attractions, and manage your travel wishlist."
	c.Data["ActivePage"]      = ""
	c.Data["PageScript"]      = ""

	 // Session — inject into every template so header.tpl can show/hide nav items
    username  := c.GetSession("username")
    firstName := c.GetSession("first_name")

    if username != nil && username != "" {
        c.Data["LoggedIn"]  = true
        c.Data["FirstName"] = firstName
    } else {
        c.Data["LoggedIn"]  = false
        c.Data["FirstName"] = ""
    }
}