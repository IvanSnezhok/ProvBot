package models

import "time"

// LogLevel represents log level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// BotLog represents a bot log entry
type BotLog struct {
	ID        int64                  `json:"id"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

