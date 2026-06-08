// controllers/home.go
// HomeController handles GET / — home page with featured countries + attractions.
package controllers

import (
	"log"

	"github.com/robiulislam99/TravelSphere/services"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Prepare() {
	c.BaseController.Prepare()
	c.Data["Title"] = "Discover Your Next Adventure"
	c.Data["ActivePage"] = "home"
}

func (c *HomeController) Get() {
	// Featured countries — degrade gracefully on API failure
	featured, err := services.Countries().GetFeatured()
	if err != nil {
		log.Printf("[HomeController] GetFeatured error: %v", err)
		featured = nil
	}

	// Popular attractions — already returns empty slice on failure
	attractions := services.Attractions().GetForHomePage()

	c.Data["FeaturedCountries"] = featured
	c.Data["PopularAttractions"] = attractions
	c.Data["PageScript"] = "home.js"

	c.Layout = "layout/base.tpl"
	c.TplName = "pages/home.tpl"
}
