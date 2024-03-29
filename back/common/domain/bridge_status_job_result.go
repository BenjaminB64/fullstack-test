package domain

import (
	"github.com/volatiletech/null"
	"time"
)

type ChabanDelmasBridgeJobResult struct {
	ID    int `json:"id"`
	JobID int `json:"job_id"`

	BoatName   string    `json:"boat_name"`
	CloseTime  time.Time `json:"close_time"`
	ReopenTime time.Time `json:"reopen_time"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}
