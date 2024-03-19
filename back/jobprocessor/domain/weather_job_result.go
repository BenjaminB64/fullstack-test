package domain

import (
	"github.com/volatiletech/null"
	"time"
)

type WeatherJobResult struct {
	ID    int
	JobID int

	Latitude         float32
	Longitude        float32
	Temperature      float32
	RelativeHumidity float32
	WeatherWmoCode   int

	CreatedAt time.Time
	UpdatedAt null.Time
	DeletedAt null.Time
}
