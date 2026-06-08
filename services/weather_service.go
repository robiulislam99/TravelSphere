// services/weather_service.go
// WeatherService fetches current weather via WeatherAPI (bonus feature).
// If the API key is missing the service silently returns nil —
// the destination page simply hides the weather card.
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// WeatherData holds the display-ready fields shown on the destination page.
type WeatherData struct {
	TempC     float64
	Condition string
	Humidity  int
	Icon      string
}

// WeatherService fetches live weather data.
type WeatherService struct {
	apiKey     string
	httpClient *http.Client
}

// NewWeatherService creates a WeatherService.
// Returns nil if WEATHER_API_KEY is not set (feature is disabled).
func NewWeatherService() *WeatherService {
	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		log.Println("[WeatherService] WEATHER_API_KEY not set — weather feature disabled")
		return nil
	}
	return &WeatherService{
		apiKey:     key,
		httpClient: &http.Client{Timeout: 8 * time.Second},
	}
}

// GetCurrent returns current weather for a city name or coordinates string.
// Returns nil, nil when the API is unavailable so callers degrade gracefully.
func (s *WeatherService) GetCurrent(query string) (*WeatherData, error) {
	if s == nil || query == "" {
		return nil, nil
	}

	url := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		s.apiKey, query,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("[WeatherService] request error: %v", err)
		return nil, nil // graceful degradation
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			Humidity  int     `json:"humidity"`
			Condition struct {
				Text string `json:"text"`
				Icon string `json:"icon"`
			} `json:"condition"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, nil
	}

	return &WeatherData{
		TempC:     raw.Current.TempC,
		Condition: raw.Current.Condition.Text,
		Humidity:  raw.Current.Humidity,
		Icon:      "https:" + raw.Current.Condition.Icon,
	}, nil
}