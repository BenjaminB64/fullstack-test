package grpc_client

import (
	"context"
	"fmt"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/common/domainProtoConverters"
	jobsproto "github.com/BenjaminB64/fullstack-test/back/common/protobuf/jobs-proto"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JobServiceClient struct {
	client jobsproto.JobServiceClient
}

func NewJobServiceClient(logger *logger.Logger, config *config.Config) (domain.JobServiceClient, error) {
	if config.JobService.Host == "" {
		return nil, fmt.Errorf("job service host is empty")
	}
	if config.JobService.Port == 0 {
		return nil, fmt.Errorf("job service port is 0")
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", config.JobService.Host, config.JobService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to connect to job service", "error", err)
		return nil, err
	}

	client := jobsproto.NewJobServiceClient(conn)

	return &JobServiceClient{client}, nil
}

func (j JobServiceClient) GetPendingJobs(ctx context.Context) ([]*commonDomain.Job, error) {
	jobs, err := j.client.GetPendingJobs(ctx, &jobsproto.GetPendingJobsRequest{})
	if err != nil {
		return nil, err
	}

	domainJobs := make([]*commonDomain.Job, len(jobs.Jobs))
	for i, job := range jobs.Jobs {
		domainJobs[i] = domainProtoConverters.GrpcJobToDomainJob(job)
	}

	return domainJobs, nil
}

// UpdateBridgeJob updates the status of a get bridge schedule job in the job service
func (j JobServiceClient) UpdateBridgeJob(ctx context.Context, jobID int, status commonDomain.JobStatus, bridgeSchedule *domain.BridgeSchedule) error {
	_, err := j.client.UpdateJob(ctx, &jobsproto.UpdateJobRequest{
		Id:     int32(jobID),
		Status: domainProtoConverters.DomainStatusToGRPCStatus(status),
		Result: &jobsproto.JobResult{
			ResultOneof: &jobsproto.JobResult_BridgeSchedule{
				BridgeSchedule: bridgeScheduleToGRPCBridgeSchedule(bridgeSchedule),
			},
		},
	})
	return err
}

// UpdateWeatherJob updates the status of a weather job in the job service
func (j JobServiceClient) UpdateWeatherJob(ctx context.Context, jobID int, status commonDomain.JobStatus, weather *domain.Weather) error {
	_, err := j.client.UpdateJob(ctx, &jobsproto.UpdateJobRequest{
		Id:     int32(jobID),
		Status: domainProtoConverters.DomainStatusToGRPCStatus(status),
		Result: &jobsproto.JobResult{
			ResultOneof: &jobsproto.JobResult_Weather{
				Weather: weatherToGRPCWeather(weather),
			},
		},
	})
	return err
}

func bridgeScheduleToGRPCBridgeSchedule(bridgeSchedule *domain.BridgeSchedule) *jobsproto.BridgeSchedule {
	bs := &jobsproto.BridgeSchedule{}
	bs.Closures = make([]*jobsproto.BridgeClosure, len(bridgeSchedule.Closures))
	for i, closure := range bridgeSchedule.Closures {
		bs.Closures[i] = &jobsproto.BridgeClosure{
			BoatName:   closure.BoatName,
			CloseTime:  timestamppb.New(closure.CloseTime),
			ReopenTime: timestamppb.New(closure.ReopenTime),
		}
	}

	return bs
}

func weatherToGRPCWeather(weather *domain.Weather) *jobsproto.Weather {
	return &jobsproto.Weather{
		RelativeHumidity: float32(weather.RelativeHumidity),
		Temperature:      int32(weather.Temperature),
		WmoCode:          int32(weather.WeatherWmoCode),
	}
}
