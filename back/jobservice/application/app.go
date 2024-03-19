package application

import (
	"context"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/database"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/http_server"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/repository"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/validator"
	"sync"
)

type Application struct {
	logger     *logger.Logger
	httpServer http_server.HTTPServerInterface
	db         *database.DB
}

func NewApplication(
	l *logger.Logger,
	c *config.Config,
	ctx context.Context,
) (*Application, error) {
	db, err := database.NewDB(ctx, l, c)
	if err != nil {
		l.Error("error creating database", "error", err)
		return nil, err
	}
	err = db.TryPing(ctx)
	if err != nil {
		l.Error("error pinging database", "error", err)
		return nil, err
	}
	err = db.EnsureSchema(ctx)
	if err != nil {
		l.Error("error ensuring schema exists", "error", err)
		return nil, err
	}

	var customValidator *validator.Validator
	customValidator, err = validator.NewValidator()
	if err != nil {
		l.Error("error creating custom validator", "error", err)
		return nil, err
	}

	jobRepository := repository.NewDBJobRepository(db)
	jobService := NewJobService(jobRepository)

	jobHandlers, _ := http_server.NewJobHandlers(l, jobService, customValidator)

	httpServer := http_server.NewHTTPServer(ctx, l, c, jobHandlers)

	return &Application{
		logger:     l,
		httpServer: httpServer,
		db:         db,
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

	err := app.httpServer.Run(ctx)
	if err != nil {
		return errors.Join(errors.New("start http server error"), err)
	}

	wg.Wait()

	return nil
}

func (app *Application) Stop() error {
	app.logger.Info("job service application is stopping")

	err := app.httpServer.Shutdown()
	if err != nil {
		return errors.Join(errors.New("stop http server error"), err)
	}

	err = app.db.Close()
	if err != nil {
		return err
	}

	return nil
}
