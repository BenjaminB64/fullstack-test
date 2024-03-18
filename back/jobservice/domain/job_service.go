package domain

import (
	"context"
)

type JobService interface {
	CreateJob(ctx context.Context, name string, taskType JobTaskType) (*Job, error)
	ReadJob(ctx context.Context, id int) (*Job, error)
	UpdateJob(ctx context.Context, job *Job) (*Job, error)
	DeleteJob(ctx context.Context, id int) error
	ListJobs(ctx context.Context) ([]*Job, error)
}
