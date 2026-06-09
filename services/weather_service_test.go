package services

import (
	"os"
	"testing"
)

// TestNewWeatherService_WithKeySet tests service creation when API key is set
func TestNewWeatherService_WithKeySet(t *testing.T) {
	// Set the API key
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "test-api-key-12345")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	service := NewWeatherService()

	if service == nil {
		t.Error("WeatherService should not be nil when API key is set")
	}
	if service.apiKey != "test-api-key-12345" {
		t.Errorf("WeatherService.apiKey = %q, want %q", service.apiKey, "test-api-key-12345")
	}
	if service.httpClient == nil {
		t.Error("WeatherService.httpClient should not be nil")
	}
}

// TestNewWeatherService_WithoutKeySet tests service creation when API key is not set
func TestNewWeatherService_WithoutKeySet(t *testing.T) {
	// Unset the API key
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Unsetenv("WEATHER_API_KEY")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	service := NewWeatherService()

	if service != nil {
		t.Error("WeatherService should be nil when API key is not set")
	}
}

// TestWeatherService_GetCurrent_NilService tests GetCurrent on nil service
func TestWeatherService_GetCurrent_NilService(t *testing.T) {
	var service *WeatherService = nil

	result, err := service.GetCurrent("London")

	if result != nil {
		t.Errorf("GetCurrent on nil service should return nil, got %v", result)
	}
	if err != nil {
		t.Errorf("GetCurrent on nil service should not return error, got %v", err)
	}
}

// TestWeatherService_GetCurrent_EmptyQuery tests GetCurrent with empty query
func TestWeatherService_GetCurrent_EmptyQuery(t *testing.T) {
	// Set the API key
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "test-api-key")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	service := NewWeatherService()

	if service == nil {
		t.Skip("Skipping test: API key not set")
	}

	result, err := service.GetCurrent("")

	// Empty query should be handled gracefully
	if result != nil {
		t.Errorf("GetCurrent with empty query should return nil, got %v", result)
	}
	if err != nil {
		t.Errorf("GetCurrent with empty query should not return error, got %v", err)
	}
}

// TestWeatherData tests the WeatherData model
func TestWeatherData_Model(t *testing.T) {
	weatherData := &WeatherData{
		TempC:     25.5,
		Condition: "Partly Cloudy",
		Humidity:  65,
		Icon:      "https://example.com/icon.png",
	}

	if weatherData.TempC != 25.5 {
		t.Errorf("TempC = %v, want 25.5", weatherData.TempC)
	}
	if weatherData.Condition != "Partly Cloudy" {
		t.Errorf("Condition = %s, want Partly Cloudy", weatherData.Condition)
	}
	if weatherData.Humidity != 65 {
		t.Errorf("Humidity = %d, want 65", weatherData.Humidity)
	}
	if weatherData.Icon != "https://example.com/icon.png" {
		t.Errorf("Icon = %s, want https://example.com/icon.png", weatherData.Icon)
	}
}

// TestWeatherService_TimeoutConfiguration tests that the service has proper timeout
func TestWeatherService_TimeoutConfiguration(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "test-key")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	service := NewWeatherService()

	if service == nil {
		t.Skip("Skipping test: could not create service")
	}

	if service.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	// Verify httpClient has a timeout (should be 8 seconds)
	if service.httpClient.Timeout == 0 {
		t.Error("httpClient should have a timeout set")
	}
}

// TestWeatherData_ZeroValues tests WeatherData with zero values
func TestWeatherData_ZeroValues(t *testing.T) {
	weatherData := &WeatherData{}

	if weatherData.TempC != 0 {
		t.Errorf("TempC = %v, want 0", weatherData.TempC)
	}
	if weatherData.Humidity != 0 {
		t.Errorf("Humidity = %d, want 0", weatherData.Humidity)
	}
	if weatherData.Condition != "" {
		t.Errorf("Condition = %s, want empty", weatherData.Condition)
	}
	if weatherData.Icon != "" {
		t.Errorf("Icon = %s, want empty", weatherData.Icon)
	}
}
