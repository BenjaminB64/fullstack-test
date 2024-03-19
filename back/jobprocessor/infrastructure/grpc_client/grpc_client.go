package grpc_client

import (
	jobsproto "github.com/BenjaminB64/fullstack-test/back/common/protobuf/jobs-proto"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JobServiceClient struct {
	client jobsproto.JobServiceClient
}

func NewJobServiceClient(logger logger.Logger, config *config.Config) (*JobServiceClient, error) {
	conn, err := grpc.Dial(
		config.JobService.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to connect to job service", "error", err)
		return nil, err
	}

	client := jobsproto.NewJobServiceClient(conn)

	return &JobServiceClient{client}, nil
}
