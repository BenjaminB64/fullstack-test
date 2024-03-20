package weather_service

import (
	"encoding/json"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"net/http"
	"time"
)

type WeatherService struct {
	httpClient *http.Client
	logger     *logger.Logger
}

// NewWeatherService creates a new WeatherService
func NewWeatherService(logger *logger.Logger) *WeatherService {
	httpClient := &http.Client{}
	httpClient.Timeout = 10 * time.Second
	return &WeatherService{
		httpClient: httpClient,
	}
}

func (ws *WeatherService) GetWeather() (*domain.Weather, error) {
	get, err := ws.httpClient.Get("https://api.open-meteo.com/v1/forecast" + "?latitude=44.8333&longitude=-0.5667&current=temperature_2m,relative_humidity_2m,weather_code")
	if err != nil {
		return nil, err
	}
	defer get.Body.Close()

	if get.StatusCode != http.StatusOK {
		ws.logger.Error("error getting weather", "status", get.StatusCode)
		return nil, errors.New("error getting weather")
	}

	var response OpenMeteoResponse
	err = json.NewDecoder(get.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	weather := &domain.Weather{
		Temperature:      response.Current.Temperature2m,
		RelativeHumidity: response.Current.RelativeHumidity2m,
		WeatherWmoCode:   response.Current.WeatherCode,
	}

	return weather, nil
}

type OpenMeteoResponse struct {
	Current struct {
		WeatherCode        int     `json:"weather_code"`
		Temperature2m      float64 `json:"temperature_2m"`
		RelativeHumidity2m float64 `json:"relative_humidity_2m"`
	} `json:"current"`
}
