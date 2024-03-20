package repository

import (
	"context"
	"database/sql"
	"errors"
	commonDomain "github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/database"
)

type DBJobRepository struct {
	db *database.DB
}

func NewDBJobRepository(db *database.DB) domain.JobRepository {
	return &DBJobRepository{db: db}
}

func (r *DBJobRepository) Create(ctx context.Context, job *commonDomain.Job) (*commonDomain.Job, error) {
	res := r.db.QueryRowContext(ctx, "INSERT INTO jobs (name, task_type, status, slack_webhook_url) VALUES ($1, $2, 'pending', $3) RETURNING id, name, status, slack_webhook_url, task_type, created_at, updated_at, deleted_at;", job.Name, job.TaskType, job.SlackWebhookURL)
	if res.Err() != nil {
		return nil, res.Err()
	}
	newJob := &commonDomain.Job{}
	err := res.Scan(&newJob.ID, &newJob.Name, &newJob.Status, &newJob.SlackWebhookURL, &newJob.TaskType, &newJob.CreatedAt, &newJob.UpdatedAt, &newJob.DeletedAt)
	if err != nil {
		return nil, err
	}

	return newJob, nil
}

func (r *DBJobRepository) Read(ctx context.Context, id int) (*commonDomain.Job, error) {
	job := &commonDomain.Job{}
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

func (r *DBJobRepository) Update(ctx context.Context, job *commonDomain.Job) error {
	err := r.db.QueryRowContext(ctx, `
		UPDATE jobs SET name = $1, task_type = $2, status = $3 WHERE id = $4 AND deleted_at IS NULL RETURNING id, name, task_type, status, slack_webhook_url, created_at, updated_at, deleted_at;
	`, job.Name, job.TaskType, job.Status, job.ID).Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.SlackWebhookURL, &job.CreatedAt, &job.UpdatedAt, &job.DeletedAt)
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

func (r *DBJobRepository) ReadLastN(ctx context.Context, n int) ([]*commonDomain.Job, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, task_type, status, slack_webhook_url, created_at, updated_at FROM jobs WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1;
	`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*commonDomain.Job, 0, n)
	for rows.Next() {
		job := &commonDomain.Job{}
		err := rows.Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.SlackWebhookURL, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *DBJobRepository) GetJobToProcess(ctx context.Context) ([]*commonDomain.Job, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, task_type, status, slack_webhook_url, created_at, updated_at FROM jobs WHERE status = 'pending' AND deleted_at IS NULL;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*commonDomain.Job, 0)
	for rows.Next() {
		job := &commonDomain.Job{}
		err := rows.Scan(&job.ID, &job.Name, &job.TaskType, &job.Status, &job.SlackWebhookURL, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *DBJobRepository) CreateWeatherJobResult(ctx context.Context, jobID int, weather *commonDomain.WeatherJobResult) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO weather_job_results (job_id, temperature, relative_humidity, weather_wmo_code, latitude, longitude) VALUES ($1, $2, $3, $4, 0, 0);
	`, jobID, weather.Temperature, weather.RelativeHumidity, weather.WeatherWmoCode)
	if err != nil {
		return errors.Join(errors.New("failed to create weather job result"), err)
	}
	return nil
}

func (r *DBJobRepository) CreateBridgeJobResult(ctx context.Context, jobID int, bridgeSchedule []*commonDomain.ChabanDelmasBridgeJobResult) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				err = errors.Join(err, txErr)
				return
			}
		}
	}()
	_, err = tx.ExecContext(ctx, `DELETE FROM chaban_delmas_bridge_job_results WHERE job_id = $1;`, jobID)
	if err != nil {
		return err
	}

	for _, closure := range bridgeSchedule {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO chaban_delmas_bridge_job_results (job_id, boat_name, closure_time, reopen_time) VALUES ($1, $2, $3, $4);
		`, jobID, closure.BoatName, closure.CloseTime, closure.ReopenTime)
		if err != nil {
			return errors.Join(errors.New("failed to create bridge closure"), err)
		}
	}

	err = tx.Commit()

	return nil
}
