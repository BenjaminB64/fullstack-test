package application

import (
	"context"
	"errors"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"time"
)

type JobProcessor struct {
	jobProcessorWorkers *JobProcessorWorkers
	jobServiceClient    domain.JobServiceClient
	jobChannel          chan *commonDomain.Job
	logger              *logger.Logger
}

func NewJobProcessor(
	jobServiceClient domain.JobServiceClient,
	logger *logger.Logger,
) domain.JobProcessor {
	jobChannel := make(chan *commonDomain.Job, 10)
	return &JobProcessor{
		jobServiceClient:    jobServiceClient,
		jobChannel:          jobChannel,
		jobProcessorWorkers: NewJobProcessorWorkers(logger, jobChannel, jobServiceClient),
		logger:              logger,
	}
}

func (j JobProcessor) ProcessJobs(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)

	workersCtx, workersCancelFn := context.WithCancel(context.Background())
	defer func() {
		workersCancelFn()
		j.jobProcessorWorkers.Wait()
		close(j.jobChannel)
	}()

	err := j.jobProcessorWorkers.Run(workersCtx)
	if err != nil {
		return errors.Join(errors.New("error running job processor workers"), err)
	}

	var jobs []*commonDomain.Job
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			jobs, err = j.jobServiceClient.GetPendingJobs(ctx)
			if err != nil {
				j.logger.Error("failed to get pending jobs", "error", err)
			}
			j.logger.Debug("got pending jobs", "numJobs", len(jobs))
			for _, job := range jobs {
				j.jobChannel <- job
			}
		}
	}
}

func (j JobProcessor) Stop() error {
	return nil
}
