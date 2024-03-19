package application

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
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

func (s *JobService) CreateJob(ctx context.Context, name string, taskType commonDomain.JobTaskType) (*commonDomain.Job, error) {
	if !taskType.IsValid() {
		return nil, domain.ErrJobInvalidTaskType
	}
	job := commonDomain.NewJob(name, taskType)

	return s.jobRepository.Create(ctx, job)
}

func (s *JobService) ReadJob(ctx context.Context, id int) (*commonDomain.Job, error) {
	return s.jobRepository.Read(ctx, id)
}

func (s *JobService) UpdateJob(ctx context.Context, job *commonDomain.Job) (*commonDomain.Job, error) {
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

func (s *JobService) ListJobs(ctx context.Context) ([]*commonDomain.Job, error) {
	return s.jobRepository.ReadLastN(ctx, 10)
}

func (s *JobService) GetJobToProcess(ctx context.Context) ([]*commonDomain.Job, error) {
	return s.jobRepository.GetJobToProcess(ctx)
}
