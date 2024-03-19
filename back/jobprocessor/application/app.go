package application

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"sync"
)

type Application struct {
	logger *logger.Logger
}

func NewApplication(
	l *logger.Logger,
	c *config.Config,
	ctx context.Context,
) (*Application, error) {

	return &Application{
		logger: l,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	app.logger.Info("job service is running")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		err := app.Stop()
		if err != nil {
			app.logger.Error("error stopping job service application", "error", err)
		}
	}()

	wg.Wait()

	return nil
}

func (app *Application) Stop() error {
	app.logger.Info("job service application is stopping")

	return nil
}
