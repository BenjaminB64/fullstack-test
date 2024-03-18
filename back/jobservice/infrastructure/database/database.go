package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/database"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DB struct {
	logger *logger.Logger
	*sql.DB
}

func NewDB(ctx context.Context, logger *logger.Logger, config *config.Config) (*DB, error) {
	pgxConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse pgx config"), err)
	}
	pgxConfig.ConnConfig.Host = config.Database.Host
	pgxConfig.ConnConfig.Port = uint16(config.Database.Port)
	pgxConfig.ConnConfig.Database = config.Database.Name
	pgxConfig.ConnConfig.User = config.Database.User
	pgxConfig.ConnConfig.Password = config.Database.Password
	pgxConfig.MaxConns = 256
	pgxConfig.MinConns = 0
	pgxConfig.HealthCheckPeriod = 10 * time.Second

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)

	if err != nil {
		return nil, errors.Join(errors.New("failed to connect to database"), err)
	}
	return &DB{
		logger: logger,
		DB:     stdlib.OpenDBFromPool(dbpool),
	}, nil
}

func (db *DB) TryPing(ctx context.Context) error {
	var err error
	for i := 0; i < 10; i++ {
		db.logger.Debug("try to connect database", "attempt", i+1)
		err = db.PingContext(ctx)
		if err == nil {
			return nil
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		db.logger.Error("failed to connect database", "error", err, "attempt", i+1)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return errors.Join(errors.New("failed to connect to the database"), err)
}

// EnsureSchema check if jobs table exists
// if not, create it (load init_db.sql)
func (db *DB) EnsureSchema(ctx context.Context) error {
	db.logger.Debug("ensure schema")
	_, err := db.ExecContext(ctx, "SELECT 1 FROM jobs;")
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "42P01" {
				_, err = db.ExecContext(ctx, database.InitDBSQL)
				if err != nil {
					return errors.Join(errors.New("failed to create table"), err)
				}
				return nil
			}
		}
		return errors.Join(errors.New("failed to check if jobs table exists"), err)
	}

	return nil
}
