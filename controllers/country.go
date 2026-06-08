// controllers/country.go
// CountryController handles:
//   GET /countries        → List()   — Country Explorer SSR page
//   GET /countries/:slug  → Detail() — Destination detail SSR page
package controllers

import (
	"log"
	"net/http"

	"github.com/robiulislam99/TravelSphere/services"
)

type CountryController struct {
	BaseController
}

func (c *CountryController) Prepare() {
	c.BaseController.Prepare()
	c.Data["ActivePage"] = "countries"
}

// List renders GET /countries with a server-side rendered country grid.
// Reads optional ?search= and ?region= query params for initial SSR filter.
func (c *CountryController) List() {
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := services.Countries().GetAll(search, region)
	if err != nil {
		log.Printf("[CountryController] GetAll error: %v", err)
		countries = nil
	}

	c.Data["Title"]        = "Explore Countries"
	c.Data["Countries"]    = countries
	c.Data["SearchQuery"]  = search
	c.Data["RegionFilter"] = region

	c.Layout  = "layout/base.tpl"
	c.TplName = "pages/countries.tpl"
}

// Detail renders GET /countries/:slug with full country info + attractions.
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")
	if slug == "" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Layout  = "layout/base.tpl"
		c.TplName = "pages/404.tpl"
		return
	}

	country, err := services.Countries().GetBySlug(slug)
	if err != nil {
		log.Printf("[CountryController] GetBySlug(%s) error: %v", slug, err)
	}

	// Unknown slug → 404
	if country == nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["Title"]  = "Not Found"
		c.Layout  = "layout/base.tpl"
		c.TplName = "pages/404.tpl"
		return
	}

	// Fetch attractions near country centre (lat/lon from model)
	attractions, _ := services.Attractions().GetByCoords(
		country.Latitude, country.Longitude, 10000, 12,
	)

	// Fetch weather (bonus) — nil when key not set
	var weather interface{}
	if ws := services.Weather(); ws != nil {
		weather, _ = ws.GetCurrent(country.Capital)
	}

	c.Data["Title"]       = country.Name
	c.Data["Country"]     = country
	c.Data["Attractions"] = attractions
	c.Data["Weather"]     = weather

	c.Layout  = "layout/base.tpl"
	c.TplName = "pages/destination.tpl"
}