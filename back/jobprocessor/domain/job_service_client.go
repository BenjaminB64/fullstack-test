package domain

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
)

type JobServiceClient interface {
	GetPendingJobs(ctx context.Context) ([]*commonDomain.Job, error)
	UpdateJobStatus(ctx context.Context, jobID int, status commonDomain.JobStatus) error
}
