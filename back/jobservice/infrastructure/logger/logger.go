package logger

import (
	"context"
	"fmt"
	"github.com/BenjaminB64/fullstack-test/back/jobservice/infrastructure/config"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(c *config.Config) (*Logger, error) {
	var h slog.Handler

	if c.App.Mode == config.AppMode_Production {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: c.App.Verbose,
		})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: c.App.Verbose,
		})
	}

	l := slog.New(h)

	return &Logger{l}, nil
}

// PrintfLogger returns a function that logs a message with the given level (used for maxprocs logger for now)
func (l *Logger) PrintfLogger(level slog.Level) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		l.Log(context.TODO(), level, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Writer(ctx context.Context, level slog.Level, msg string) io.Writer {
	return writeFunc(func(data []byte) (int, error) {
		l.Log(ctx, level, msg, "data", string(data))
		return 0, nil
	})
}

// writeFunc convert func to io.Writer
type writeFunc func([]byte) (int, error)

func (fn writeFunc) Write(data []byte) (int, error) {
	return fn(data)
}
