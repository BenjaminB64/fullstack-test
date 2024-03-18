package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/database"
)

type DBJobRepository struct {
	db *database.DB
}

func NewDBJobRepository(db *database.DB) domain.JobRepository {
	return &DBJobRepository{db: db}
}

func (r *DBJobRepository) Create(ctx context.Context, job *domain.Job) (*domain.Job, error) {
	res := r.db.QueryRowContext(ctx, "INSERT INTO jobs (name, task_type, status) VALUES ($1, $2, 'pending') RETURNING *;", job.Name, job.TaskType)
	if res.Err() != nil {
		return nil, res.Err()
	}
	newJob := &domain.Job{}
	err := res.Scan(&newJob.ID, &newJob.Name, &newJob.Status, &newJob.TaskType, &newJob.CreatedAt, &newJob.UpdatedAt, &newJob.DeletedAt)
	if err != nil {
		return nil, err
	}

	return newJob, nil
}

func (r *DBJobRepository) Read(ctx context.Context, id int) (*domain.Job, error) {
	job := &domain.Job{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, task_type, status, created_at, updated_at FROM jobs WHERE id = $1 AND deleted_at IS NULL;
	`, id).Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrJobNotFound
		}
		return nil, err
	}
	return job, nil
}

func (r *DBJobRepository) Update(ctx context.Context, job *domain.Job) error {
	err := r.db.QueryRowContext(ctx, `
		UPDATE jobs SET name = $1, task_type = $2, status = $3 WHERE id = $4 AND deleted_at IS NULL RETURNING *;
	`, job.Name, job.TaskType, job.Status, job.ID).Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrJobNotFound
		}
		return err
	}
	return nil
}

func (r *DBJobRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE jobs SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBJobRepository) ReadLastN(ctx context.Context, n int) ([]*domain.Job, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, task_type, status, created_at, updated_at FROM jobs WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1;
	`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*domain.Job, 0, n)
	for rows.Next() {
		job := &domain.Job{}
		err := rows.Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
