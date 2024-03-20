package http_server

import (
	"errors"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/http_server/dtos"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator2 "github.com/go-playground/validator/v10"
	"github.com/volatiletech/null"
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

// CreateJob godoc
//
//	@Summary		Create a job
//	@Description	Create a job
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			createJobRequest	body		dtos.CreateJobRequest	true	"Create Job Request"
//	@Success		201					{object}	dtos.JobResponse
//	@Failure		400					{object}	dtos.ApiError
//	@Failure		500					{object}	dtos.ApiError
//	@Router			/jobs [post]
func (h *JobHandlers) CreateJob(c *gin.Context) {
	createJobRequest := dtos.CreateJobRequest{}
	err := c.ShouldBindJSON(&createJobRequest)

	if err != nil {
		h.logger.Debug("create job validation error", "error", err)
		apiError := &dtos.ApiError{
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
	job := &commonDomain.Job{}

	job.Name = createJobRequest.Name
	job.TaskType = commonDomain.JobTaskType(createJobRequest.TaskType)

	if createJobRequest.SlackWebhook != "" {
		job.SlackWebhookURL = null.StringFrom(createJobRequest.SlackWebhook)
	}

	job, err = h.jobService.CreateJob(c, job)
	if err != nil {
		h.logger.Error("failed to create job", "error", err)
		c.JSON(http.StatusInternalServerError, &dtos.ApiError{
			Error: "Internal error, failed to create job, please retry later",
		})
		return
	}

	res := &dtos.JobResponse{}
	res.FromDomain(job)

	c.JSON(http.StatusCreated, res)
}

// GetJobs godoc
//
//	@Summary		Get jobs
//	@Description	Get jobs
//	@Tags			jobs
//	@Produce		json
//	@Success		200	{object}	dtos.JobsResponse
//	@Failure		500	{object}	dtos.ApiError
//	@Router			/jobs [get]
func (h *JobHandlers) GetJobs(c *gin.Context) {
	jobs, err := h.jobService.ListJobs(c)
	if err != nil {
		h.logger.Error("failed to list jobs", "error", err)
		c.JSON(http.StatusInternalServerError, &dtos.ApiError{
			Error: "Internal error, failed to list jobs, please retry later",
		})
		return
	}

	res := dtos.JobsResponse{}
	res.FromDomain(jobs)

	c.JSON(http.StatusOK, res)
}

// GetJob godoc
//
//	@Summary		Get a job by ID
//	@Description	Get a job by ID
//	@Tags			jobs
//	@Produce		json
//	@Param			id	path		int	true	"Job ID"
//	@Success		200	{object}	dtos.JobResponse
//	@Failure		400	{object}	dtos.ApiError
//	@Failure		404	{object}	dtos.ApiError
//	@Failure		500	{object}	dtos.ApiError
//	@Router			/jobs/{id} [get]
func (h *JobHandlers) GetJob(c *gin.Context) {
	var getJobURI dtos.GetJobURI

	err := c.ShouldBindUri(&getJobURI)
	if err != nil {
		h.logger.Debug("get job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &dtos.ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	job, err := h.jobService.ReadJob(c, getJobURI.ID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &dtos.ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to read job", "error", err)
		c.JSON(http.StatusInternalServerError, &dtos.ApiError{
			Error: "Internal error, failed to read job, please retry later",
		})
		return
	}

	res := &dtos.JobResponse{}
	res.FromDomain(job)

	c.JSON(http.StatusOK, res)
}

// UpdateJob godoc
//
//	@Summary		Update a job by ID
//	@Description	Update a job by ID
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			id					path	int						true	"Job ID"
//	@Param			updateJobRequest	body	dtos.UpdateJobRequest	true	"Update Job Request"
//	@Success		204
//	@Failure		400	{object}	dtos.ApiError
//	@Failure		404	{object}	dtos.ApiError
//	@Failure		500	{object}	dtos.ApiError
//	@Router			/jobs/{id} [put]
func (h *JobHandlers) UpdateJob(c *gin.Context) {
	var updateJobURI dtos.UpdateJobURI

	err := c.ShouldBindUri(&updateJobURI)
	if err != nil {
		h.logger.Debug("update job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &dtos.ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	updateJobRequest := dtos.UpdateJobRequest{}
	err = c.ShouldBindJSON(&updateJobRequest)
	if err != nil {
		h.logger.Debug("update job validation error", "error", err)
		apiError := &dtos.ApiError{
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
	job := &commonDomain.Job{}

	job.ID = updateJobURI.ID
	job.Name = updateJobRequest.Name
	job.TaskType = commonDomain.JobTaskType(updateJobRequest.TaskType)
	job.Status = commonDomain.JobStatus(updateJobRequest.Status)

	_, err = h.jobService.UpdateJob(c, job)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &dtos.ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to update job", "error", err)
		c.JSON(http.StatusInternalServerError, &dtos.ApiError{
			Error: "Internal error, failed to update job, please retry later",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteJob godoc
//
//	@Summary		Delete a job by ID
//	@Description	Delete a job by ID
//	@Tags			jobs
//	@Param			id	path	int	true	"Job ID"
//	@Success		204
//	@Failure		400	{object}	dtos.ApiError
//	@Failure		404	{object}	dtos.ApiError
//	@Failure		500	{object}	dtos.ApiError
//	@Router			/jobs/{id} [delete]
func (h *JobHandlers) DeleteJob(c *gin.Context) {
	var deleteJobURI dtos.DeleteJobURI

	err := c.ShouldBindUri(&deleteJobURI)
	if err != nil {
		h.logger.Debug("delete job URI bind error", "error", err)
		c.JSON(http.StatusBadRequest, &dtos.ApiError{
			Error: "Invalid job ID",
		})
		return
	}

	err = h.jobService.DeleteJob(c, deleteJobURI.ID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, &dtos.ApiError{
				Error: "Job not found",
			})
			return
		}
		h.logger.Error("failed to delete job", "error", err)
		c.JSON(http.StatusInternalServerError, &dtos.ApiError{
			Error: "Internal error, failed to delete job, please retry later",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
