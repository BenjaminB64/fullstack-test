package domain

import (
	"github.com/volatiletech/null"
	"time"
)

type Job struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Status   JobStatus   `json:"status"`
	TaskType JobTaskType `json:"task_type"`

	Weather      *WeatherJobResult             `json:"weather"`
	BridgeStatus []ChabanDelmasBridgeJobResult `json:"bridge_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

func NewJob(name string, jobType JobTaskType) *Job {
	return &Job{
		Name:     name,
		Status:   JobStatus_Pending,
		TaskType: jobType,
	}
}
