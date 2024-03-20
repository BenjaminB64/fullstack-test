package main

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/bridge_service"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/grpc_client"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/weather_service"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/application"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	var err error

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	var c *config.Config
	c, err = config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	var l *logger.Logger
	l, err = logger.NewLogger(c)
	if err != nil {
		log.Fatal(err)
	}

	l.Debug("show config", "config", c)

	_, err = maxprocs.Set(maxprocs.Logger(l.PrintfLogger(slog.LevelDebug)))
	if err != nil {
		l.Error("error setting maxprocs", "error", err)
		return
	}

	jobServiceClient, err := grpc_client.NewJobServiceClient(l, c)
	if err != nil {
		l.Error("error creating job service client", "error", err)
		return
	}
	weatherService := weather_service.NewWeatherService(l)
	bridgeService := bridge_service.NewBridgeService(l)

	jobProcessor := application.NewJobProcessor(jobServiceClient, l, weatherService, bridgeService)

	var app *application.Application
	app, err = application.NewApplication(l, c, ctx, jobProcessor)
	if err != nil {
		l.Error("error creating application", "error", err)
		return
	}

	err = app.Run(ctx)
	if err != nil {
		l.Error("error running application", "error", err)
		return
	}

	l.Info("application stopped")

}
