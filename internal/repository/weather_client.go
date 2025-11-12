package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"weather-api-lab/internal/domain"
	"weather-api-lab/internal/dto"
)

type WeatherClient interface {
	GetTemperatureByLocation(location *domain.Location) (float64, error)
}

type weatherClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewWeatherClient(baseURL, apiKey string) WeatherClient {
	return &weatherClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTemperatureByLocation busca a temperatura pela localização
func (c *weatherClient) GetTemperatureByLocation(location *domain.Location) (float64, error) {
	if location == nil || location.City == "" {
		return 0, domain.ErrInvalidLocation
	}

	query := fmt.Sprintf("%s, %s, Brazil", location.City, location.State)
	requestURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", c.baseURL, c.apiKey, url.QueryEscape(query))

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return 0, fmt.Errorf("error querying WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return 0, domain.ErrWeatherNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error in WeatherAPI: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading WeatherAPI response: %w", err)
	}

	var weatherResp dto.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return 0, fmt.Errorf("error parsing WeatherAPI response: %w", err)
	}

	return weatherResp.Current.TempC, nil
}
