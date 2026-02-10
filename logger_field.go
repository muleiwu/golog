package golog

import (
	"github.com/muleiwu/gsr"
)

// LoggerField represents a key-value pair for structured logging
type LoggerField struct {
	Key   string
	Value any
}

// Field creates a new LoggerField with the given key and value
// This is the primary way to add structured context to log messages
//
// Example:
//
//	logger.Info("user logged in", Field("user_id", 123), Field("ip", "192.168.1.1"))
func Field(key string, value any) gsr.LoggerField {
	return &LoggerField{
		Key:   key,
		Value: value,
	}
}

// GetKey returns the field's key
func (f *LoggerField) GetKey() string {
	return f.Key
}

// GetValue returns the field's value
func (f *LoggerField) GetValue() any {
	return f.Value
}
