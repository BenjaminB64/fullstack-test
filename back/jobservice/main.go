package main

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/BenjaminB64/fullstack-test/back/jobservice/application"
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

	var app *application.Application
	app, err = application.NewApplication(l, c, ctx)
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
