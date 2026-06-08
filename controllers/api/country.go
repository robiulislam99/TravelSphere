// controllers/api/country.go
// CountryAPIController — JSON endpoints consumed by AJAX on /countries.
package api

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/services"
	"github.com/robiulislam99/TravelSphere/utils"
)

type CountryAPIController struct {
	web.Controller
}

// GetAll handles GET /api/countries?search=...&region=...
func (c *CountryAPIController) GetAll() {
	search := c.GetString("search")
	region := c.GetString("region")

	if region != "" && !utils.IsValidRegion(region) {
		utils.SendBadRequest(&c.Controller, "invalid region value")
		return
	}

	countries, err := services.Countries().GetAll(search, region)
	if err != nil {
		utils.SendInternalError(&c.Controller)
		return
	}
	utils.SendSuccess(&c.Controller, countries)
}

// GetBySlug handles GET /api/countries/:slug
func (c *CountryAPIController) GetBySlug() {
	slug := c.Ctx.Input.Param(":slug")
	if !utils.IsValidSlug(slug) {
		utils.SendBadRequest(&c.Controller, "invalid slug format")
		return
	}

	country, err := services.Countries().GetBySlug(slug)
	if err != nil {
		utils.SendInternalError(&c.Controller)
		return
	}
	if country == nil {
		utils.SendNotFound(&c.Controller, "country not found")
		return
	}
	utils.SendSuccess(&c.Controller, country)
}