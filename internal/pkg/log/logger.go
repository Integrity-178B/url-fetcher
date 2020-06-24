package log

import (
	"log"
	"os"
)

// Logger encapsulates basic logger
type Logger struct {
	*log.Logger
}

// NewLogger creates new logger instance with specified prefix
func NewLogger(prefix string) *Logger {
	return &Logger{log.New(os.Stdout, prefix, log.LstdFlags)}
}
