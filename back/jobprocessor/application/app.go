package application

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
)

type Application struct {
	logger       *logger.Logger
	jobProcessor domain.JobProcessor
}

func NewApplication(
	l *logger.Logger,
	c *config.Config,
	ctx context.Context,
	jobProcessor domain.JobProcessor,
) (*Application, error) {

	return &Application{
		logger:       l,
		jobProcessor: jobProcessor,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	app.logger.Info("job service is running")

	err := app.jobProcessor.ProcessJobs(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) Stop() error {
	app.logger.Info("job service application is stopping")

	return nil
}
