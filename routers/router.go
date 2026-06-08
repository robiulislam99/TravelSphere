package routers

import (
	"github.com/robiulislam99/TravelSphere/controllers"
	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/filters"
)

func init() {
	// ── Logging filter — applies to all routes ────────────────
	web.InsertFilter("/*", web.BeforeRouter, filters.LoggingFilter)
 
	// ── Auth filter — protects wishlist and dashboard ─────────
	web.InsertFilter("/wishlist", web.BeforeRouter, filters.AuthFilter)
	web.InsertFilter("/dashboard", web.BeforeRouter, filters.AuthFilter)
 
	// ── SSR page routes ───────────────────────────────────────
	web.Router("/", &controllers.HomeController{})
	web.Router("/countries", &controllers.CountryController{}, "get:List")
	web.Router("/countries/:slug", &controllers.CountryController{}, "get:Detail")
	web.Router("/wishlist", &controllers.WishlistController{})
	web.Router("/dashboard", &controllers.DashboardController{})
}