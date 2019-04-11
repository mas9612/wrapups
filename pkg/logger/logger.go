package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger represents the simple logger used in elastic package.
type Logger struct {
	Logger *zap.Logger
}

// Printf prints given message with zap package.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Debug("debug log", zap.String("msg", fmt.Sprintf(format, v...)))
}
