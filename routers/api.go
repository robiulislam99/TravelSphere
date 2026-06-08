// routers/api.go
// Registers all /api/* JSON routes.

package routers

import (
	"github.com/beego/beego/v2/server/web"
	apiControllers "github.com/robiulislam99/TravelSphere/controllers/api"
	"github.com/robiulislam99/TravelSphere/filters"
)

func init() {
	// Auth filter protects wishlist and dashboard API endpoints
	web.InsertFilter("/api/wishlist", web.BeforeRouter, filters.AuthFilter)
	web.InsertFilter("/api/wishlist/*", web.BeforeRouter, filters.AuthFilter)
	web.InsertFilter("/api/dashboard/*", web.BeforeRouter, filters.AuthFilter)

	// ── Country API ───────────────────────────────────────────
	// GET /api/countries?search=...&region=...
	web.Router("/api/countries", &apiControllers.CountryAPIController{}, "get:GetAll")
	// GET /api/countries/:slug
	web.Router("/api/countries/:slug", &apiControllers.CountryAPIController{}, "get:GetBySlug")

	// ── Attractions API ───────────────────────────────────────
	// GET /api/attractions?lat=...&lon=...
	web.Router("/api/attractions", &apiControllers.AttractionAPIController{}, "get:GetByCoords")

	// ── Wishlist API (CRUD) ───────────────────────────────────
	// GET    /api/wishlist
	// POST   /api/wishlist
	// PUT    /api/wishlist/:id
	// DELETE /api/wishlist/:id
	web.Router("/api/wishlist", &apiControllers.WishlistAPIController{}, "get:GetAll;post:Create")
	web.Router("/api/wishlist/:id", &apiControllers.WishlistAPIController{}, "put:Update;delete:Delete")

	// ── Dashboard API ─────────────────────────────────────────
	// GET /api/dashboard/summary
	web.Router("/api/dashboard/summary", &apiControllers.DashboardAPIController{}, "get:Summary")
}