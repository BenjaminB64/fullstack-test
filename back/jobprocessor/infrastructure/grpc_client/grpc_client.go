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

func (j JobServiceClient) UpdateJobStatus(ctx context.Context, jobID int, status commonDomain.JobStatus) error {
	_, err := j.client.UpdateJobStatus(ctx, &jobsproto.UpdateJobStatusRequest{
		Job: &jobsproto.Job{
			Id:     int32(jobID),
			Status: domainProtoConverters.DomainStatusToGRPCStatus(status),
		},
	})
	return err
}
