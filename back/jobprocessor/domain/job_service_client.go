package domain

import commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"

type JobServiceClient interface {
	GetPendingJobs() ([]*commonDomain.Job, error)
	UpdateJobStatus(jobID string, status commonDomain.JobStatus) error
}
