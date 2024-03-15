package application

import (
	"context"

	"github.com/BenjaminB64/fullstack-test/back/jobservice/pkg/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/pkg/logger"
)

type JobService struct {
	logger *logger.Logger
}

func NewJobService(
	logger *logger.Logger,
	config *config.Config,
	ctx context.Context,
) (*JobService, error) {
	return &JobService{
		logger: logger,
	}, nil
}

func (js *JobService) Run() {
	js.logger.Info("job service is running")
}
