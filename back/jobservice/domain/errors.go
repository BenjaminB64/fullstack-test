package domain

import "errors"

var ErrJobNotFound = errors.New("job not found")
var ErrJobInvalidTaskType = errors.New("job task type is invalid")
var ErrJobInvalidStatus = errors.New("job status is invalid")
