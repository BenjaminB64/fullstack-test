package config

import (
	"errors"
	"log/slog"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Database struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		Name     string `koanf:"name"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
	} `koanf:"database"`
	App struct {
		Port    int        `koanf:"port"`
		Verbose slog.Level `koanf:"verbose"`
		Mode    AppMode    `koanf:"mode"`
	} `koanf:"app"`
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
