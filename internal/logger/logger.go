package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "" {
		envLogLevel = zerolog.LevelInfoValue
	}

	logLevel, err := zerolog.ParseLevel(envLogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05 01-02"}).
		With().
		Timestamp().
		Logger()

	return logger
}
