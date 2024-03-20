package application

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/database"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/repository"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"testing"
	"time"
)

type JobServiceTestSuite struct {
	suite.Suite
	JobService domain.JobService
	db         *database.DB
}

func (s *JobServiceTestSuite) SetupSuite() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	l := &logger.Logger{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}

	c, err := config.NewConfig()
	s.Require().NoError(err)
	s.Require().NotNil(c)

	db, err := database.NewDB(ctx, l, c)
	s.Require().NoError(err)
	s.Require().NotNil(db)

	err = db.TryPing(ctx)
	s.Require().NoError(err)

	_, err = db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
	s.Require().NoError(err)

	err = db.EnsureSchema(ctx)
	s.Require().NoError(err)

	jobRepository := repository.NewDBJobRepository(db)
	s.Require().NotNil(jobRepository)

	s.db = db

	s.JobService = NewJobService(jobRepository)
	s.Require().NotNil(s.JobService)

}

func (s *JobServiceTestSuite) TestCreateJob() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: commonDomain.JobTaskType_Weather,
	}

	job, err := s.JobService.CreateJob(ctx, newJob)
	s.NoError(err)
	s.NotNil(job)
	s.NotEqual(0, job.ID)
	s.Equal("test", job.Name)
	s.Equal(commonDomain.JobStatus_Pending, job.Status)
	s.Equal(commonDomain.JobTaskType_Weather, job.TaskType)
	s.NotNil(job.CreatedAt)
	s.Nil(job.UpdatedAt.Ptr())
	s.Nil(job.DeletedAt.Ptr())
}

func (s *JobServiceTestSuite) TestCreateJobInvalidTaskType() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: "invalid_type",
	}

	job, err := s.JobService.CreateJob(ctx, newJob)
	s.ErrorIs(err, domain.ErrJobInvalidTaskType)
	s.Nil(job)
}

func (s *JobServiceTestSuite) TestReadJob() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: commonDomain.JobTaskType_Weather,
	}

	job, err := s.JobService.CreateJob(ctx, newJob)
	s.NoError(err)
	s.NotNil(job)

	readJob, err := s.JobService.ReadJob(ctx, job.ID)
	s.NoError(err)
	s.NotNil(readJob)
	s.Equal(job.ID, readJob.ID)
	s.Equal(job.Name, readJob.Name)
	s.Equal(job.Status, readJob.Status)
	s.Equal(job.TaskType, readJob.TaskType)
	s.Equal(job.CreatedAt, readJob.CreatedAt)
	s.Equal(job.UpdatedAt, readJob.UpdatedAt)
	s.Equal(job.DeletedAt, readJob.DeletedAt)
}

func (s *JobServiceTestSuite) TestReadJobNotFound() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	readJob, err := s.JobService.ReadJob(ctx, 147)

	s.ErrorIs(err, domain.ErrJobNotFound)
	s.Nil(readJob)
}

func (s *JobServiceTestSuite) TestUpdateJob() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: commonDomain.JobTaskType_Weather,
	}

	job, err := s.JobService.CreateJob(ctx, newJob)
	s.NoError(err)
	s.NotNil(job)

	job.Status = commonDomain.JobStatus_Completed
	_, err = s.JobService.UpdateJob(ctx, job)
	s.NoError(err)

	readJob, err := s.JobService.ReadJob(ctx, job.ID)
	s.NoError(err)
	s.NotNil(readJob)
	s.Equal(commonDomain.JobStatus_Completed, readJob.Status)

}

func (s *JobServiceTestSuite) TestUpdateJobInvalidStatus() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: commonDomain.JobTaskType_Weather,
	}
	job, err := s.JobService.CreateJob(ctx, newJob)
	s.NoError(err)
	s.NotNil(job)

	job.Status = "invalid_status"
	_, err = s.JobService.UpdateJob(ctx, job)
	s.ErrorIs(err, domain.ErrJobInvalidStatus)
}

func (s *JobServiceTestSuite) TearDownSuite() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	_, err := s.db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
	s.NoError(err)

	err = s.db.Close()
	s.NoError(err)
}

func (s *JobServiceTestSuite) TestDeleteJob() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	newJob := &commonDomain.Job{
		Name:     "test",
		TaskType: commonDomain.JobTaskType_Weather,
	}

	job, err := s.JobService.CreateJob(ctx, newJob)
	s.NoError(err)
	s.NotNil(job)

	err = s.JobService.DeleteJob(ctx, job.ID)
	s.NoError(err)

	readJob, err := s.JobService.ReadJob(ctx, job.ID)

	s.ErrorIs(err, domain.ErrJobNotFound)
	s.Nil(readJob)
}

func TestJobServiceSuite(t *testing.T) {
	suite.Run(t, new(JobServiceTestSuite))
}
