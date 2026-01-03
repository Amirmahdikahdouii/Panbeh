package logger

import (
	"log/slog"
	"os"

	"github.com/panbeh/otp-backend/internal/config"
)

func New(level config.LogLevel) *slog.Logger {
	var l slog.Level
	switch level {
	case config.LogLevelDebug:
		l = slog.LevelDebug
	case config.LogLevelInfo:
		l = slog.LevelInfo
	case config.LogLevelWarn:
		l = slog.LevelWarn
	case config.LogLevelError:
		l = slog.LevelError
	default:
		l = slog.LevelDebug
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l})
	return slog.New(h)
}
