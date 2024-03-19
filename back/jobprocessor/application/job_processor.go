package application

import (
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
)

type JobProcessor struct {
	jobServiceClient domain.JobServiceClient
}

func NewJobProcessor(jobServiceClient domain.JobServiceClient) domain.JobProcessor {
	return &JobProcessor{
		jobServiceClient: jobServiceClient,
	}
}

func (j JobProcessor) ProcessJobs() error {
	//TODO implement me
	panic("implement me")
}
