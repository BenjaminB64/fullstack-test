package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/BenjaminB64/fullstack-test/back/jobservice/application"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/pkg/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/pkg/logger"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.NewLogger(config)
	if err != nil {
		log.Fatal(err)
	}

	maxprocs.Set(maxprocs.Logger(logger.PrintfLogger(slog.LevelDebug)))

	logger.Debug("show config", "config", config)

	ctx := context.Background()

	app, err := application.NewJobService(logger, config, ctx)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
