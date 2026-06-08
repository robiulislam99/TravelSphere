// controllers/home.go
// HomeController handles GET / — the home page.
// Renders home.tpl with featured countries and popular attractions.
// Full implementation comes in Phase 5.
package controllers

// HomeController handles the home page SSR route.
type HomeController struct {
	BaseController
}

// Prepare sets page-level template data before Get() runs.
func (c *HomeController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"] = "Home"
	c.Data["ActivePage"] = "home"
}

// Get renders the home page with featured countries and attractions.
// TODO (Phase 5): inject real data from CountryService and AttractionService.
func (c *HomeController) Get() {
	c.Data["FeaturedCountries"] = []interface{}{}
	c.Data["PopularAttractions"] = []interface{}{}

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/home.tpl"
}