package application

import (
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
)

type JobService struct {
	jobRepository domain.JobRepository
}

func NewJobService(jobRepository domain.JobRepository) domain.JobService {
	return &JobService{
		jobRepository: jobRepository,
	}
}
