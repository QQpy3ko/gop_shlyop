package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(env string) zerolog.Logger {
	var l zerolog.Logger

	switch env {
	case envLocal:
		// local - to console
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger()
	case envDev:
		// dev/prod - JSON
		l = zerolog.New(os.Stdout).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger()
	case envProd:
		l = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	default:
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger()
	}

	return l
}
