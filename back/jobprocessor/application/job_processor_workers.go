package application

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"sync"
)

type JobProcessorWorkers struct {
	logger *logger.Logger
	ctx    context.Context

	wg         *sync.WaitGroup
	jobChannel <-chan *commonDomain.Job

	jobServiceClient domain.JobServiceClient

	weatherService domain.WeatherService
	bridgeService  domain.BridgeService
}

func NewJobProcessorWorkers(
	logger *logger.Logger,
	jobChannel <-chan *commonDomain.Job,
	jobServiceClient domain.JobServiceClient,
	weatherService domain.WeatherService,
	bridgeService domain.BridgeService,
) *JobProcessorWorkers {
	return &JobProcessorWorkers{
		wg:               &sync.WaitGroup{},
		jobChannel:       jobChannel,
		logger:           logger,
		jobServiceClient: jobServiceClient,
		weatherService:   weatherService,
		bridgeService:    bridgeService,
	}
}

func (j *JobProcessorWorkers) Run(ctx context.Context) error {
	j.ctx = ctx
	for i := 0; i < 2; i++ {
		j.AddWorker()
	}

	return nil
}

func (j *JobProcessorWorkers) AddWorker() {
	j.wg.Add(1)
	go j.RunWorker()
}

func (j *JobProcessorWorkers) RunWorker() {
	defer j.wg.Done()
	for {
		select {
		case <-j.ctx.Done():
			return
		case job := <-j.jobChannel:
			// process job
			j.logger.Info("processing job", "job", job)
			err := j.ProcessJob(j.ctx, job)
			if err != nil {
				j.logger.Error("failed to process job", "error", err)
			}
		}
	}
}

func (j *JobProcessorWorkers) ProcessJob(ctx context.Context, job *commonDomain.Job) error {
	switch job.TaskType {
	case commonDomain.JobTaskType_Weather:
		weather, err := j.weatherService.GetWeather()
		if err != nil {
			return err
		}
		j.logger.Debug("got weather", "weather", weather)
		err = j.jobServiceClient.UpdateWeatherJob(ctx, job.ID, commonDomain.JobStatus_Completed, weather)
		if err != nil {
			return err
		}
	case commonDomain.JobTaskType_GetChabanDelmasBridgeSchedule:
		bridgeSchedule, err := j.bridgeService.GetBridgeSchedule()
		if err != nil {
			return err
		}
		j.logger.Debug("got bridge schedule", "bridgeSchedule", bridgeSchedule)
		err = j.jobServiceClient.UpdateBridgeJob(ctx, job.ID, commonDomain.JobStatus_Completed, bridgeSchedule)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JobProcessorWorkers) Wait() {
	j.wg.Wait()
	j.logger.Info("all workers have stopped")
}
