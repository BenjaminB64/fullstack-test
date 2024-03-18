package http_server

import (
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/http_server/dtos"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator2 "github.com/go-playground/validator/v10"
	"net/http"
)

type JobHandlers struct {
	logger          *logger.Logger
	jobService      domain.JobService
	customValidator *validator.Validator
}

func NewJobHandlers(logger *logger.Logger, jobService domain.JobService, customValidator *validator.Validator) (*JobHandlers, error) {
	err := customValidator.RegisterOn(binding.Validator.Engine().(*validator2.Validate))
	if err != nil {
		return nil, errors.Join(errors.New("failed to register custom validator"), err)
	}

	return &JobHandlers{
		logger:          logger,
		jobService:      jobService,
		customValidator: customValidator,
	}, nil
}

func (h *JobHandlers) CreateJob(c *gin.Context) {
	createJobRequest := dtos.CreateJobRequest{}
	err := c.ShouldBindJSON(&createJobRequest)

	if err != nil {
		h.logger.Debug("create job validation error", "error", err)
		apiError := &ApiError{
			Error: err.Error(),
		}
		var validationErrors validator2.ValidationErrors
		if errors.As(err, &validationErrors) {
			apiError.Fields, err = h.customValidator.GetValidationErrorsMap(validationErrors)
			if err != nil {
				h.logger.Error("failed to get validation errors map", "error", err)
				c.JSON(http.StatusInternalServerError, apiError)
				return
			}
		}
		c.JSON(http.StatusBadRequest, apiError)
		return
	}
	var job *domain.Job
	job, err = h.jobService.CreateJob(c, createJobRequest.Name, domain.JobTaskType(createJobRequest.TaskType))
	if err != nil {
		h.logger.Error("failed to create job", "error", err)
		c.JSON(http.StatusInternalServerError, &ApiError{
			Error: "Internal error, failed to create job, please retry later",
		})
		return
	}

	c.JSON(http.StatusCreated, job)
}

func (h *JobHandlers) GetJobs(c *gin.Context) {
	jobs, err := h.jobService.ListJobs(c)
	if err != nil {
		h.logger.Error("failed to list jobs", "error", err)
		c.JSON(http.StatusInternalServerError, &ApiError{
			Error: "Internal error, failed to list jobs, please retry later",
		})
		return
	}

	res := dtos.JobsResponse{}
	res.FromDomain(jobs)

	c.JSON(http.StatusOK, res)
}

func (h *JobHandlers) GetJob(c *gin.Context) {
	var getJobURI dtos.GetJobURI

	err := c.ShouldBindUri(&getJobURI)
	if err != nil {
		h.logger.Debug("get job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	job, err := h.jobService.ReadJob(c, getJobURI.ID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to read job", "error", err)
		c.JSON(http.StatusInternalServerError, &ApiError{
			Error: "Internal error, failed to read job, please retry later",
		})
		return
	}

	res := &dtos.JobResponse{}
	res.FromDomain(job)

	c.JSON(http.StatusOK, res)
}

func (h *JobHandlers) UpdateJob(c *gin.Context) {
	var updateJobURI dtos.UpdateJobURI

	err := c.ShouldBindUri(&updateJobURI)
	if err != nil {
		h.logger.Debug("update job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	updateJobRequest := dtos.UpdateJobRequest{}
	err = c.ShouldBindJSON(&updateJobRequest)
	if err != nil {
		h.logger.Debug("update job validation error", "error", err)
		apiError := &ApiError{
			Error: err.Error(),
		}
		var validationErrors validator2.ValidationErrors
		if errors.As(err, &validationErrors) {
			apiError.Fields, err = h.customValidator.GetValidationErrorsMap(validationErrors)
			if err != nil {
				h.logger.Error("failed to get validation errors map", "error", err)
				c.JSON(http.StatusInternalServerError, apiError)
				return
			}
		}
		c.JSON(http.StatusBadRequest, apiError)
		return
	}
	job := &domain.Job{}

	job.ID = updateJobURI.ID
	job.Name = updateJobRequest.Name
	job.TaskType = updateJobRequest.TaskType
	job.Status = updateJobRequest.Status

	_, err = h.jobService.UpdateJob(c, job)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to update job", "error", err)
		c.JSON(http.StatusInternalServerError, &ApiError{
			Error: "Internal error, failed to update job, please retry later",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *JobHandlers) DeleteJob(c *gin.Context) {
	var deleteJobURI dtos.DeleteJobURI

	err := c.ShouldBindUri(&deleteJobURI)
	if err != nil {
		h.logger.Debug("delete job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	err = h.jobService.DeleteJob(c, deleteJobURI.ID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to delete job", "error", err)
		c.JSON(http.StatusInternalServerError, &ApiError{
			Error: "Internal error, failed to delete job, please retry later",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

type ApiError struct {
	Error  string                          `json:"error"`
	Fields map[string]validator.FieldError `json:"fields,omitempty"`
}
