package config

import (
	"log/slog"
	"os"
)

func InitLog() {
	opts := &slog.HandlerOptions{
		Level:     getLoggerLevel(),
		AddSource: true,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func getLoggerLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
