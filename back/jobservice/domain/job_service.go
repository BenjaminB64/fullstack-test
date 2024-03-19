package domain

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
)

type JobService interface {
	CreateJob(ctx context.Context, name string, taskType commonDomain.JobTaskType) (*commonDomain.Job, error)
	ReadJob(ctx context.Context, id int) (*commonDomain.Job, error)
	UpdateJob(ctx context.Context, job *commonDomain.Job) (*commonDomain.Job, error)
	DeleteJob(ctx context.Context, id int) error
	ListJobs(ctx context.Context) ([]*commonDomain.Job, error)
	GetJobToProcess(ctx context.Context) ([]*commonDomain.Job, error)
}
