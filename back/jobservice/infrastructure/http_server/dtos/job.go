package dtos

import (
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/volatiletech/null"
	"time"
)

// data transfer objects

type CreateJobRequest struct {
	TaskType     TaskType `json:"taskType" binding:"required,enum"`
	Name         string   `json:"name" binding:"required,max=255"`
	SlackWebhook string   `json:"slackWebhook" binding:"slack_webhook_url"`
}

type JobResponse struct {
	ID int `json:"id"`

	Name     string   `json:"name"`
	TaskType TaskType `json:"taskType"`
	Status   Status   `json:"status"`

	SlackWebhook null.String `json:"slackWebhook" swaggertype:"string" example:"https://hooks.slack.com/services/..."`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt null.Time `json:"updatedAt" swaggertype:"string" example:"2021-01-01T00:00:00Z"`
}

type UpdateJobRequest struct {
	// in: body
	Name     string   `json:"name" validate:"required,max=255"`
	TaskType TaskType `json:"taskType" validate:"required,enum"`
	Status   Status   `json:"status" validate:"required,enum"`
}

type TaskType = domain.JobTaskType
type Status = domain.JobStatus

func (jr *JobResponse) FromDomain(job *domain.Job) {
	if job == nil || jr == nil {
		return
	}
	*jr = JobResponse{
		ID: job.ID,

		Name:     job.Name,
		TaskType: job.TaskType,
		Status:   job.Status,

		SlackWebhook: null.String{},

		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

type JobsResponse []*JobResponse

func (jr *JobsResponse) FromDomain(jobs []*domain.Job) {
	if jr == nil {
		return
	}
	*jr = make(JobsResponse, 0, len(jobs))
	for _, job := range jobs {
		j := &JobResponse{}
		j.FromDomain(job)
		*jr = append(*jr, j)
	}
}

type IdJobURI struct {
	ID int `uri:"id" validate:"required"`
}

type GetJobURI = IdJobURI
type DeleteJobURI = IdJobURI
type UpdateJobURI = IdJobURI
