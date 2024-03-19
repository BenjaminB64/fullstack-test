package domain

import "context"

type JobProcessor interface {
	ProcessJobs(ctx context.Context) error
}
