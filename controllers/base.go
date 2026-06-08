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
}