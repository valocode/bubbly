package env

import (
	"os"

	"github.com/rs/zerolog"
)

// NewDebugLogger is a convenience function for generating a new
// zerolog.Logger in debug mode for use in bubbly tests
func NewDebugLogger() *zerolog.Logger {
	// Initialize Logger

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	}).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	return &logger
}
