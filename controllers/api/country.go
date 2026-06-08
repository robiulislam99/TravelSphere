// controllers/api/country.go
// CountryAPIController handles JSON API routes for country data.
// Full implementation in Phase 6.
package api

import "github.com/beego/beego/v2/server/web"

// CountryAPIController serves JSON responses for country data.
type CountryAPIController struct {
	web.Controller
}

// GetAll handles GET /api/countries?search=...&region=...
func (c *CountryAPIController) GetAll() {
	c.Data["json"] = map[string]interface{}{"data": []interface{}{}, "message": "stub — Phase 6"}
	c.ServeJSON()
}

// GetBySlug handles GET /api/countries/:slug
func (c *CountryAPIController) GetBySlug() {
	c.Data["json"] = map[string]interface{}{"data": nil, "message": "stub — Phase 6"}
	c.ServeJSON()
}