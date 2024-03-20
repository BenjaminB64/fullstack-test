package domain

type WeatherService interface {
	GetWeather() (*Weather, error)
}

// value object
type Weather struct {
	Temperature      float64
	RelativeHumidity float64
	WeatherWmoCode   int
}
