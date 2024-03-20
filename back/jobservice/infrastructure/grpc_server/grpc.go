package grpc_server

import (
	"context"
	"errors"
	"fmt"
	domain2 "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/common/domainProtoConverters"
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

func (s *GrpcServer) GetPendingJobs(ctx context.Context, request *jobsproto.GetPendingJobsRequest) (*jobsproto.GetPendingJobsResponse, error) {
	jobs, err := s.jobService.GetJobToProcess(ctx)
	if err != nil {
		return nil, err
	}
	res := &jobsproto.GetPendingJobsResponse{}
	for _, job := range jobs {
		res.Jobs = append(res.Jobs, domainProtoConverters.DomainJobToGRPCJob(job))
	}

	return res, nil
}

func (s *GrpcServer) UpdateJob(ctx context.Context, request *jobsproto.UpdateJobRequest) (*jobsproto.UpdateJobResponse, error) {
	job, err := s.jobService.ReadJob(ctx, int(request.Id))
	if err != nil {
		return nil, err
	}

	job.Status = domainProtoConverters.GrpcStatusToDomainStatus(request.Status)

	if job.TaskType == domain2.JobTaskType_GetChabanDelmasBridgeSchedule {
		if bridgeSchedule, ok := request.Result.ResultOneof.(*jobsproto.JobResult_BridgeSchedule); ok {
			bridgeResult := domainProtoConverters.GrpcBridgeScheduleToBridgeJobResult(bridgeSchedule.BridgeSchedule)
			err = s.jobService.CreateBridgeJobResult(ctx, job.ID, bridgeResult)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("job type is GetChabanDelmasBridgeSchedule but result is not BridgeSchedule")
		}
	}

	if job.TaskType == domain2.JobTaskType_Weather {
		if weather, ok := request.Result.ResultOneof.(*jobsproto.JobResult_Weather); ok {
			weatherResult := domainProtoConverters.GrpcWeatherToWeatherJobResult(weather.Weather)
			err = s.jobService.CreateWeatherJobResult(ctx, job.ID, weatherResult)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("job type is Weather but result is not Weather")
		}
	}

	_, err = s.jobService.UpdateJob(ctx, job)
	if err != nil {
		return nil, err
	}

	return &jobsproto.UpdateJobResponse{}, nil
}
