package dtos

import (
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
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
	Name     string   `json:"name" validate:"required,max=255"`
	TaskType TaskType `json:"taskType" validate:"required,enum"`
	Status   Status   `json:"status" validate:"required,enum"`
}

// FromDomain converts a domain Job to a JobResponse
// we assume that taskType and status are same as in commonDomain
func (jr *JobResponse) FromDomain(job *commonDomain.Job) {
	if job == nil || jr == nil {
		return
	}
	*jr = JobResponse{
		ID: job.ID,

		Name:     job.Name,
		TaskType: TaskType(job.TaskType),
		Status:   Status(job.Status),

		SlackWebhook: job.SlackWebhookURL,

		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

type JobsResponse []*JobResponse

func (jr *JobsResponse) FromDomain(jobs []*commonDomain.Job) {
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

/*
swag cli doesn't support workspace, so it can't handle the following types
type TaskType = commonDomain.JobTaskType
type Status = commonDomain.JobStatus
*/

type TaskType string

const (
	TaskType_Unknown                       TaskType = ""
	TaskType_Weather                       TaskType = "get_weather"
	TaskType_GetChabanDelmasBridgeSchedule TaskType = "get_chaban_delmas_bridge_schedule"
)

func (j TaskType) IsValid() bool {
	switch j {
	case TaskType_Weather, TaskType_GetChabanDelmasBridgeSchedule:
		return true
	}
	return false
}

type Status string

const (
	Status_Unknown    Status = ""
	Status_Pending    Status = "pending"
	Status_InProgress Status = "in_progress"
	Status_Completed  Status = "completed"
	Status_Failed     Status = "failed"
)

func (j Status) IsValid() bool {
	switch j {
	case Status_Pending, Status_InProgress, Status_Completed, Status_Failed:
		return true
	}
	return false
}
