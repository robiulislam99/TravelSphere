// controllers/base.go
// BaseController is embedded by all SSR page controllers.
// Its Prepare() method sets common template data available to every page:
//   - Title, MetaDescription
//   - ActivePage (for nav highlighting)
//   - PageScript (optional per-page JS file name)
package controllers

import "github.com/beego/beego/v2/server/web"

// BaseController embeds web.Controller and provides shared Prepare logic.
type BaseController struct {
	web.Controller
}

// Prepare runs before every action in any controller that embeds BaseController.
// Set default template data here so every page template has these variables.
func (c *BaseController) Prepare() {
	// Default page title — individual controllers should override this
	c.Data["Title"] = "TravelSphere"
	c.Data["MetaDescription"] = "Discover destinations, explore attractions, and manage your travel wishlist."

	// ActivePage is used by the header partial to highlight the current nav link.
	// Override in each controller's action method.
	c.Data["ActivePage"] = ""

	// PageScript is the filename (without path) of a page-specific JS file.
	// The base layout appends /static/js/<PageScript> if set.
	c.Data["PageScript"] = ""
	c.Layout = "layouts/base.tpl"
}