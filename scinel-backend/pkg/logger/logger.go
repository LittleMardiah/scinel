package logger

import (
"log/slog"
"os"
)

// Logger wraps slog.Logger for structured logging.
type Logger struct {
*slog.Logger
}

// NewLogger creates a new logger with JSON format
func NewLogger(level string) *Logger {
var lvl slog.Level
switch level {
case "debug":
lvl = slog.LevelDebug
case "warn":
lvl = slog.LevelWarn
case "error":
lvl = slog.LevelError
default:
lvl = slog.LevelInfo
}

handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
Level: lvl,
})
return &Logger{Logger: slog.New(handler)}
}