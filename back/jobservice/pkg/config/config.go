package config

import (
	"errors"
	"log/slog"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Database struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	} `koanf:"database"`
	App struct {
		Port    int        `koanf:"port"`
		Verbose slog.Level `koanf:"verbose"`
	} `koanf:"app"`
}

func NewConfig() (*Config, error) {
	k := koanf.New(".")
	var config Config

	// Load config from file
	if err := k.Load(env.Provider("", "_", nil), nil); err != nil {
		return nil, errors.New("failed to load config file: " + err.Error())
	}

	k.Unmarshal("", &config)

	return &config, nil
}
