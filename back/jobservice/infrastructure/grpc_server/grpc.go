package grpc_server

import (
	"context"
	"errors"
	"fmt"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	jobsproto "github.com/BenjaminB64/fullstack-test/back/common/protobuf/jobs-proto"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	*jobsproto.UnimplementedJobServiceServer
	server     *grpc.Server
	jobService domain.JobService
	logger     *logger.Logger
	port       int
}

var _ jobsproto.JobServiceServer = &GrpcServer{}

func NewGrpcServer(logger *logger.Logger, config *config.Config, jobService domain.JobService) *GrpcServer {
	return &GrpcServer{
		logger:     logger,
		port:       config.GRPCServer.Port,
		jobService: jobService,
	}
}

func (s *GrpcServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	s.server = grpc.NewServer(opts...)
	jobsproto.RegisterJobServiceServer(s.server, s)
	s.logger.Info("grpc server started", "port", s.port)

	err = s.server.Serve(lis)
	if err != nil {
		return errors.Join(errors.New("failed to serve grpc"), err)
	}

	return nil
}

func (s *GrpcServer) Stop() error {
	if s.server == nil {
		return nil
	}
	s.server.GracefulStop()
	s.logger.Info("grpc server stopped", "port", s.port)
	return nil
}

func (s *GrpcServer) GetJobToProcess(ctx context.Context, request *jobsproto.GetPendingJobsRequest) (*jobsproto.GetPendingJobsResponse, error) {
	jobs, err := s.jobService.GetJobToProcess(ctx)
	if err != nil {
		return nil, err
	}
	res := &jobsproto.GetPendingJobsResponse{}
	for _, job := range jobs {
		res.Jobs = append(res.Jobs, domainJobToGRPCJob(job))
	}

	return res, nil
}

func (s *GrpcServer) UpdateJobStatus(ctx context.Context, request *jobsproto.UpdateJobStatusRequest) (*jobsproto.UpdateJobStatusResponse, error) {
	job, err := s.jobService.ReadJob(ctx, int(request.Job.Id))
	if err != nil {
		return nil, err
	}

	job.Status = grpcStatusToDomainStatus(request.Job.Status)
	_, err = s.jobService.UpdateJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return &jobsproto.UpdateJobStatusResponse{}, nil
}

func domainJobToGRPCJob(job *commonDomain.Job) *jobsproto.Job {
	return &jobsproto.Job{
		Id:     int32(job.ID),
		Type:   domainTaskTypeToGRPCTaskType(job.TaskType),
		Status: domainStatusToGRPCStatus(job.Status),
		Result: nil,
	}
}

func domainTaskTypeToGRPCTaskType(taskType commonDomain.JobTaskType) jobsproto.JobType {
	switch taskType {
	case commonDomain.JobTaskType_Weather:
		return jobsproto.JobType_JOB_TYPE_GET_WEATHER
	case commonDomain.JobTaskType_GetChabanDelmasBridgeStatus:
		return jobsproto.JobType_JOB_TYPE_GET_CHABAN_DELMAS_BRIDGE_STATUS
	default:
		return jobsproto.JobType_JOB_TYPE_UNKNOWN
	}
}

func domainStatusToGRPCStatus(status commonDomain.JobStatus) jobsproto.JobStatus {
	switch status {
	case commonDomain.JobStatus_Pending:
		return jobsproto.JobStatus_JOB_STATUS_PENDING
	case commonDomain.JobStatus_InProgress:
		return jobsproto.JobStatus_JOB_STATUS_IN_PROGRESS
	case commonDomain.JobStatus_Completed:
		return jobsproto.JobStatus_JOB_STATUS_COMPLETED
	case commonDomain.JobStatus_Failed:
		return jobsproto.JobStatus_JOB_STATUS_FAILED
	default:
		return jobsproto.JobStatus_JOB_STATUS_UNKNOWN
	}
}

func grpcJobToDomainJob(job *jobsproto.Job) *commonDomain.Job {
	return &commonDomain.Job{
		ID:     int(job.Id),
		Status: grpcStatusToDomainStatus(job.Status),
	}
}

func grpcStatusToDomainStatus(status jobsproto.JobStatus) commonDomain.JobStatus {
	switch status {
	case jobsproto.JobStatus_JOB_STATUS_PENDING:
		return commonDomain.JobStatus_Pending
	case jobsproto.JobStatus_JOB_STATUS_IN_PROGRESS:
		return commonDomain.JobStatus_InProgress
	case jobsproto.JobStatus_JOB_STATUS_COMPLETED:
		return commonDomain.JobStatus_Completed
	case jobsproto.JobStatus_JOB_STATUS_FAILED:
		return commonDomain.JobStatus_Failed
	default:
		return commonDomain.JobStatus_Unknown
	}
}
