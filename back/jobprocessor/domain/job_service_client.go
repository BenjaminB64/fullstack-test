package domain

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
)

type JobServiceClient interface {
	GetPendingJobs(ctx context.Context) ([]*commonDomain.Job, error)

	UpdateBridgeJob(ctx context.Context, jobID int, status commonDomain.JobStatus, bridgeSchedule *BridgeSchedule) error
	UpdateWeatherJob(ctx context.Context, jobID int, status commonDomain.JobStatus, weather *Weather) error
}
