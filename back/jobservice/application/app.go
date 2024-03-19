package application

import (
	"context"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/database"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/grpc_server"
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
	grpcServer *grpc_server.GrpcServer
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
	grpcServer := grpc_server.NewGrpcServer(l, c, jobService)

	return &Application{
		logger:     l,
		httpServer: httpServer,
		db:         db,
		grpcServer: grpcServer,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	app.logger.Info("job service is running")

	var cancelFn context.CancelCauseFunc
	ctx, cancelFn = context.WithCancelCause(ctx)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.RunGRPCServer(ctx)
		if err != nil {
			app.logger.Error("error running grpc server", "error", err)
			cancelFn(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.RunHTTPServer(ctx)
		if err != nil {
			app.logger.Error("error running http server", "error", err)
			cancelFn(err)
		}
	}()

	<-ctx.Done()
	err := app.Stop()
	if err != nil {
		if ctx.Err() != nil {
			return errors.Join(ctx.Err(), err)
		}
		return err
	}
	wg.Wait()

	if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	return nil
}

func (app *Application) RunHTTPServer(ctx context.Context) error {
	err := app.httpServer.Run(ctx)
	if err != nil {
		return errors.Join(errors.New("start http server error"), err)
	}
	return nil
}

func (app *Application) RunGRPCServer(ctx context.Context) error {
	app.logger.Info("job service grpc server is running")

	err := app.grpcServer.Run()
	if err != nil {
		return errors.Join(errors.New("start grpc server error"), err)
	}

	return nil
}

func (app *Application) Stop() error {
	app.logger.Info("job service application is stopping")

	err := app.httpServer.Shutdown()
	if err != nil {
		return errors.Join(errors.New("stop http server error"), err)
	}

	err = app.grpcServer.Stop()
	if err != nil {
		return errors.Join(errors.New("stop grpc server error"), err)
	}

	err = app.db.Close()
	if err != nil {
		return err
	}

	return nil
}
