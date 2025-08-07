package logger

import (
	"os"
	"todo_list/config"

	"golang.org/x/exp/slog"
)

var LogLevel = map[string]slog.Level{
	"DEBUG":    slog.LevelDebug,
	"INFO":     slog.LevelInfo,
	"WARNING":  slog.LevelWarn,
	"ERROR":    slog.LevelError,
	"CRITICAL": slog.LevelError,
	"FATAL":    slog.LevelError,
}

func New(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     LogLevel[cfg.LogLevel],
		AddSource: true,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
