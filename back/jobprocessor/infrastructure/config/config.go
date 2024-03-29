package config

import (
	"errors"
	"log/slog"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App struct {
		Verbose slog.Level `koanf:"verbose"`
		Mode    AppMode    `koanf:"mode"`
	} `koanf:"app"`
	JobService struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	} `koanf:"jobservice"`
}

type AppMode string

const (
	AppMode_Development AppMode = "development"
	AppMode_Production  AppMode = "production"
	AppMode_Test        AppMode = "test"
)

func NewConfig() (*Config, error) {
	k := koanf.New(".")
	var config Config

	// Load config from environment variables
	if err := k.Load(env.Provider("", "_", nil), nil); err != nil {
		return nil, errors.Join(errors.New("failed to load env"), err)
	}

	err := k.Unmarshal("", &config)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal config"), err)
	}

	return &config, nil
}
