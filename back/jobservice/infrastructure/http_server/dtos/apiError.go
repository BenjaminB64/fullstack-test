package dtos

import "github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/validator"

type ApiError struct {
	Error  string                          `json:"error"`
	Fields map[string]validator.FieldError `json:"fields,omitempty"`
}
