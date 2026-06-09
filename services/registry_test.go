package services

import (
	"sync"
	"testing"
)

// TestInit_FirstCall tests that Init initializes all services on first call
func TestInit_FirstCall(t *testing.T) {
	// Save original values
	originalOnce := once

	// Reset the once to allow re-initialization for testing
	once = sync.Once{}

	Init()

	if countries == nil {
		t.Error("countries service should be initialized")
	}
	if attractions == nil {
		t.Error("attractions service should be initialized")
	}
	if wishlist == nil {
		t.Error("wishlist service should be initialized")
	}
	if dashboard == nil {
		t.Error("dashboard service should be initialized")
	}
	// weather may be nil if WEATHER_API_KEY is not set, which is expected

	// Restore original state
	once = originalOnce
}

// TestInit_Idempotent tests that Init can be called multiple times safely
func TestInit_Idempotent(t *testing.T) {
	// Save original values
	originalOnce := once
	originalCountries := countries
	originalAttractions := attractions
	originalWishlist := wishlist
	originalDashboard := dashboard

	// Reset the once
	once = sync.Once{}
	countries = nil
	attractions = nil
	wishlist = nil
	dashboard = nil

	// First call
	Init()
	firstCountries := countries
	firstAttractions := attractions
	firstWishlist := wishlist
	firstDashboard := dashboard

	// Second call
	Init()
	secondCountries := countries
	secondAttractions := attractions
	secondWishlist := wishlist
	secondDashboard := dashboard

	// Services should be the same instances (not recreated)
	if firstCountries != secondCountries {
		t.Error("countries service changed after second Init call")
	}
	if firstAttractions != secondAttractions {
		t.Error("attractions service changed after second Init call")
	}
	if firstWishlist != secondWishlist {
		t.Error("wishlist service changed after second Init call")
	}
	if firstDashboard != secondDashboard {
		t.Error("dashboard service changed after second Init call")
	}

	// Restore original state
	once = originalOnce
	countries = originalCountries
	attractions = originalAttractions
	wishlist = originalWishlist
	dashboard = originalDashboard
}

// TestCountries_ReturnsInstance tests that Countries() returns the singleton
func TestCountries_ReturnsInstance(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalCountries := countries

	once = sync.Once{}
	countries = nil

	Init()

	service1 := Countries()
	service2 := Countries()

	if service1 == nil {
		t.Error("Countries() returned nil")
	}
	if service1 != service2 {
		t.Error("Countries() should return the same instance")
	}

	// Restore
	once = originalOnce
	countries = originalCountries
}

// TestAttractions_ReturnsInstance tests that Attractions() returns the singleton
func TestAttractions_ReturnsInstance(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalAttractions := attractions

	once = sync.Once{}
	attractions = nil

	Init()

	service1 := Attractions()
	service2 := Attractions()

	if service1 == nil {
		t.Error("Attractions() returned nil")
	}
	if service1 != service2 {
		t.Error("Attractions() should return the same instance")
	}

	// Restore
	once = originalOnce
	attractions = originalAttractions
}

// TestWishlist_ReturnsInstance tests that Wishlist() returns the singleton
func TestWishlist_ReturnsInstance(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalWishlist := wishlist

	once = sync.Once{}
	wishlist = nil

	Init()

	service1 := Wishlist()
	service2 := Wishlist()

	if service1 == nil {
		t.Error("Wishlist() returned nil")
	}
	if service1 != service2 {
		t.Error("Wishlist() should return the same instance")
	}

	// Restore
	once = originalOnce
	wishlist = originalWishlist
}

// TestDashboard_ReturnsInstance tests that Dashboard() returns the singleton
func TestDashboard_ReturnsInstance(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalDashboard := dashboard

	once = sync.Once{}
	dashboard = nil

	Init()

	service1 := Dashboard()
	service2 := Dashboard()

	if service1 == nil {
		t.Error("Dashboard() returned nil")
	}
	if service1 != service2 {
		t.Error("Dashboard() should return the same instance")
	}

	// Restore
	once = originalOnce
	dashboard = originalDashboard
}

// TestWeather_ReturnsInstance tests that Weather() returns the singleton
func TestWeather_ReturnsInstance(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalWeather := weather

	once = sync.Once{}
	weather = nil

	Init()

	service1 := Weather()
	service2 := Weather()

	// Weather may be nil if API key is not set, which is valid
	if service1 != service2 {
		t.Error("Weather() should return the same instance (even if nil)")
	}

	// Restore
	once = originalOnce
	weather = originalWeather
}

// TestRegistry_DependencyWiring tests that services are properly wired
func TestRegistry_DependencyWiring(t *testing.T) {
	// Save and reset state
	originalOnce := once
	originalDashboard := dashboard
	originalWishlist := wishlist

	once = sync.Once{}
	dashboard = nil
	wishlist = nil

	Init()

	dashboardService := Dashboard()
	wishlistService := Wishlist()

	if dashboardService == nil {
		t.Error("DashboardService should be initialized")
	}
	if wishlistService == nil {
		t.Error("WishlistService should be initialized")
	}

	// Dashboard should have a reference to the Wishlist service
	if dashboardService.wishlist != wishlistService {
		t.Error("Dashboard should be wired with the Wishlist service")
	}

	// Restore
	once = originalOnce
	dashboard = originalDashboard
	wishlist = originalWishlist
}
