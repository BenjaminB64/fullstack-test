package application

import (
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"testing"
	"time"
)

type JobProcessorTestSuite struct {
	suite.Suite
	JobProcessor domain.JobProcessor
}

func (s *JobProcessorTestSuite) SetupSuite() {
	_, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	_ = &logger.Logger{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}

	c, err := config.NewConfig()
	s.Require().NoError(err)
	s.Require().NotNil(c)

	s.JobProcessor = NewJobProcessor(nil)
	s.Require().NotNil(s.JobProcessor)

}

func TestJobProcessorSuite(t *testing.T) {
	suite.Run(t, new(JobProcessorTestSuite))
}
