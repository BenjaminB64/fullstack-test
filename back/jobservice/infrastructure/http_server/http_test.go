package http_server

import (
	"bytes"
	"context"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/validator"
	domain2 "github.com/BenjaminB64/fullstack-test/back/jobservice/mocks/github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type HttpTestSuite struct {
	suite.Suite
	HTTPServer     *HTTPServer
	jobServiceMock *domain2.MockJobService
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTestSuite))
}

func (s *HttpTestSuite) SetupSuite() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFn()

	l := &logger.Logger{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}

	c, err := config.NewConfig()
	c.App.Mode = config.AppMode_Test
	s.Require().NoError(err)
	s.Require().NotNil(c)

	customValidator, err := validator.NewValidator()
	s.Require().NoError(err)

	jobService := domain2.NewMockJobService(s.T())
	jobHandlers, err := NewJobHandlers(l, jobService, customValidator)
	s.Require().NoError(err)

	s.jobServiceMock = jobService

	s.HTTPServer = NewHTTPServer(
		ctx,
		l,
		c,
		jobHandlers,
	).(*HTTPServer)
}

func (s *HttpTestSuite) TestCreateJob() {
	s.jobServiceMock.Mock.On("CreateJob", mock.Anything, "Test Job", domain.JobTaskType_Weather).Return(&domain.Job{
		ID:           1,
		Name:         "Test Job",
		Status:       domain.JobStatus_Pending,
		TaskType:     domain.JobTaskType_Weather,
		Weather:      nil,
		BridgeStatus: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    null.Time{},
		DeletedAt:    null.Time{},
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer([]byte(`{"name":"Test Job","taskType":"get_weather"}`)))
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusCreated, resp.Code)
}

func (s *HttpTestSuite) TestCreateJobInvalid() {
	req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer([]byte(`{"name":"Test Job"}`)))
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusBadRequest, resp.Code)
}

func (s *HttpTestSuite) TestGetJobs() {
	s.jobServiceMock.Mock.On("ListJobs", mock.Anything).Return([]*domain.Job{
		{
			ID:           1,
			Name:         "Test Job",
			Status:       domain.JobStatus_Pending,
			TaskType:     domain.JobTaskType_Weather,
			Weather:      nil,
			BridgeStatus: nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    null.Time{},
			DeletedAt:    null.Time{},
		},
		{
			ID:           2,
			Name:         "Test Job 2",
			Status:       domain.JobStatus_Completed,
			TaskType:     domain.JobTaskType_GetChabanDelmasBridgeStatus,
			Weather:      nil,
			BridgeStatus: nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    null.Time{},
			DeletedAt:    null.Time{},
		},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusOK, resp.Code)
}

func (s *HttpTestSuite) TestGetJob() {
	s.jobServiceMock.Mock.On("ReadJob", mock.Anything, 1).Return(&domain.Job{
		ID:           1,
		Name:         "Test Job",
		Status:       domain.JobStatus_Pending,
		TaskType:     domain.JobTaskType_Weather,
		Weather:      nil,
		BridgeStatus: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    null.Time{},
		DeletedAt:    null.Time{},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/jobs/1", nil)
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusOK, resp.Code)
}

func (s *HttpTestSuite) TestUpdateJob() {
	s.jobServiceMock.Mock.On("UpdateJob", mock.Anything, mock.Anything).Return(&domain.Job{}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/jobs/1", bytes.NewBuffer([]byte(`{"name":"Updated Job","taskType":"get_weather","status":"pending"}`)))
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusNoContent, resp.Code)
}

func (s *HttpTestSuite) TestDeleteJob() {
	s.jobServiceMock.Mock.On("DeleteJob", mock.Anything, 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/jobs/1", nil)
	resp := httptest.NewRecorder()

	s.HTTPServer.server.Handler.ServeHTTP(resp, req)

	s.Equal(http.StatusNoContent, resp.Code)
}
