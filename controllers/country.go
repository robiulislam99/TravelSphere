// controllers/country.go
// CountryController handles SSR routes:
//   GET /countries        → List() renders countries.tpl
//   GET /countries/:slug  → Detail() renders destination.tpl
//
// Full implementation with real service calls comes in Phase 5.
package controllers

// CountryController handles the Country Explorer and Destination Detail pages.
type CountryController struct {
	BaseController
}

// Prepare sets shared template data for country pages.
func (c *CountryController) Prepare() {
	c.BaseController.Prepare()
	c.Data["ActivePage"] = "countries"
}

// List renders GET /countries — the Country Explorer page.
// TODO (Phase 5): call CountryService.SearchCountries and pass results.
func (c *CountryController) List() {
	c.Data["Title"] = "Explore Countries"
	c.Data["Countries"] = []interface{}{}
	c.Data["SearchQuery"] = c.GetString("search")
	c.Data["RegionFilter"] = c.GetString("region")

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/countries.tpl"
}

// Detail renders GET /countries/:slug — the Destination Detail page.
// TODO (Phase 5): parse slug, call CountryService.GetBySlug + AttractionService.
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")
	if slug == "" {
		c.Abort("404")
		return
	}

	c.Data["Title"] = slug
	c.Data["Country"] = map[string]interface{}{
		"Name":                slug,
		"Slug":                slug,
		"FlagURL":             "",
		"Capital":             "—",
		"Region":              "—",
		"Subregion":           "—",
		"FormattedPopulation": "—",
		"CurrencyDisplay":     "—",
		"LanguageDisplay":     "—",
	}
	c.Data["Attractions"] = []interface{}{}
	c.Data["Weather"] = nil

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/destination.tpl"
}