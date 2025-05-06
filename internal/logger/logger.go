package logger

import (
	"fiber-gorm/internal/config"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configures the global logger based on application configuration
func Setup(cfg config.Config) {
	// Set global log level based on config
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Use pretty console logging in development
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// In production, use structured JSON logging
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Log startup information
	log.Info().
		Str("environment", cfg.Environment).
		Str("log_level", level.String()).
		Msg("Logger initialized")
}

// Logger returns a new logger with the given component name
func New(component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}

// Error logs an error with context
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal error and exits
func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Debug().Msg(fmt.Sprintf(msg, args...))
	} else {
		log.Debug().Msg(msg)
	}
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Info().Msg(fmt.Sprintf(msg, args...))
	} else {
		log.Info().Msg(msg)
	}
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Warn().Msg(fmt.Sprintf(msg, args...))
	} else {
		log.Warn().Msg(msg)
	}
}
