package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	config2 "github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
	"time"
)

type HTTPServer struct {
	router gin.IRouter
	logger *logger.Logger
	server *http.Server
}

type HTTPServerInterface interface {
	Run(context.Context) error
	Shutdown() error
}

//	@title			Jobs API Service
//	@version		1.0
//	@description	This service provides a RESTful API for managing jobs that can be executed asynchronously
//	@contact.name	BenjaminB64
//	@host			localhost:8080
//	@BasePath		/

func NewHTTPServer(
	ctx context.Context,
	logger *logger.Logger,
	config *config.Config,
	jobHandlers *JobHandlers,
) HTTPServerInterface {
	httpServer := &HTTPServer{
		logger: logger,
	}

	r := gin.New()
	if config.App.Mode == config2.AppMode_Production {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.App.Mode == config2.AppMode_Test {
		gin.SetMode(gin.TestMode)
	}

	gin.DefaultWriter = logger.Writer(ctx, slog.LevelInfo, "gin")
	gin.DefaultErrorWriter = logger.Writer(ctx, slog.LevelError, "gin")

	r.Use(sloggin.NewWithConfig(logger.Logger, sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

		Filters: []sloggin.Filter{},
	}))

	r.Use(gin.Recovery())

	r.HandleMethodNotAllowed = false

	r.POST("/jobs", jobHandlers.CreateJob)
	r.GET("/jobs", jobHandlers.GetJobs)
	r.GET("/jobs/:id", jobHandlers.GetJob)
	r.PUT("/jobs/:id", jobHandlers.UpdateJob)
	r.DELETE("/jobs/:id", jobHandlers.DeleteJob)

	httpServer.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", config.App.Port),
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       50 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return httpServer
}

func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Info("http server listening", "bindAddr", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		// graceful shutdown
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

func (s *HTTPServer) Shutdown() error {
	s.logger.Info("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.logger.Info("http server is down")

	return nil
}
