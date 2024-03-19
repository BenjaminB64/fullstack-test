package domain

import (
	"context"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job *commonDomain.Job) (*commonDomain.Job, error)
	Read(ctx context.Context, id int) (*commonDomain.Job, error)
	Update(ctx context.Context, job *commonDomain.Job) error
	Delete(ctx context.Context, id int) error
	ReadLastN(ctx context.Context, n int) ([]*commonDomain.Job, error)
	GetJobToProcess(ctx context.Context) ([]*commonDomain.Job, error)
}
