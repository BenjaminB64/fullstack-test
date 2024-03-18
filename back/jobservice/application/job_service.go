package application

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
)

type JobService struct {
	jobRepository domain.JobRepository
}

func NewJobService(jobRepository domain.JobRepository) domain.JobService {
	return &JobService{
		jobRepository: jobRepository,
	}
}

func (s *JobService) CreateJob(ctx context.Context, name string, taskType domain.JobTaskType) (*domain.Job, error) {
	if !taskType.IsValid() {
		return nil, domain.ErrJobInvalidTaskType
	}
	job := domain.NewJob(name, taskType)

	return s.jobRepository.Create(ctx, job)
}

func (s *JobService) ReadJob(ctx context.Context, id int) (*domain.Job, error) {
	return s.jobRepository.Read(ctx, id)
}

func (s *JobService) UpdateJob(ctx context.Context, job *domain.Job) (*domain.Job, error) {
	if !job.TaskType.IsValid() {
		return nil, domain.ErrJobInvalidTaskType
	}
	if !job.Status.IsValid() {
		return nil, domain.ErrJobInvalidStatus
	}

	return nil, s.jobRepository.Update(ctx, job)
}

func (s *JobService) DeleteJob(ctx context.Context, id int) error {
	return s.jobRepository.Delete(ctx, id)
}

func (s *JobService) ListJobs(ctx context.Context) ([]*domain.Job, error) {
	return s.jobRepository.ReadLastN(ctx, 10)
}
