// services/registry.go
// Registry instantiates all services once at startup and provides
// package-level accessors so controllers can obtain services without
// passing dependencies through constructors.
//
// This avoids global vars scattered across files while keeping wiring simple
// for an assessment-scale project.
package services

import "sync"

var (
	once      sync.Once
	countries *CountryService
	attractions *AttractionService
	wishlist  *WishlistService
	dashboard *DashboardService
	weather   *WeatherService
)

// Init must be called once from main.go (or init in routers) before any
// request is handled. It is safe to call multiple times — only the first call
// has effect thanks to sync.Once.
func Init() {
	once.Do(func() {
		wishlist    = NewWishlistService()
		countries   = NewCountryService()
		attractions = NewAttractionService()
		dashboard   = NewDashboardService(wishlist)
		weather     = NewWeatherService() // nil if key not set
	})
}

// Countries returns the singleton CountryService.
func Countries() *CountryService { return countries }

// Attractions returns the singleton AttractionService.
func Attractions() *AttractionService { return attractions }

// Wishlist returns the singleton WishlistService.
func Wishlist() *WishlistService { return wishlist }

// Dashboard returns the singleton DashboardService.
func Dashboard() *DashboardService { return dashboard }

// Weather returns the singleton WeatherService (may be nil).
func Weather() *WeatherService { return weather }