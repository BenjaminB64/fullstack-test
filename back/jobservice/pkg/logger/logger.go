package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/BenjaminB64/fullstack-test/back/jobservice/pkg/config"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(config *config.Config) (*Logger, error) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.App.Verbose,
	})

	l := slog.New(h)

	return &Logger{l}, nil
}

func (l *Logger) PrintfLogger(level slog.Level) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		l.Log(context.TODO(), level, fmt.Sprintf(format, args...))
	}
}
