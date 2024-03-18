package domain

import (
	"context"
)

type JobRepository interface {
	Create(ctx context.Context, job *Job) (*Job, error)
	Read(ctx context.Context, id int) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id int) error
	ReadLastN(ctx context.Context, n int) ([]*Job, error)
}
