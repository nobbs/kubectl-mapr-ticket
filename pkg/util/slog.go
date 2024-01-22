package util

import (
	"io"
	"log/slog"
	"time"

	"github.com/charmbracelet/log"
)

// SetupLogging sets up logging for the application, to be used in the root command
// preRun hook.
func SetupLogging(out io.Writer, withDebug bool) error {
	level := log.InfoLevel
	if withDebug {
		level = log.DebugLevel
	}

	handler := log.NewWithOptions(out, log.Options{
		Level:           level,
		ReportTimestamp: true,
		TimeFormat:      time.Stamp,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
