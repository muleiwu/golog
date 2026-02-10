package main

import (
	"github.com/muleiwu/golog"
)

func main() {
	// Create a development logger
	logger, err := golog.NewDevelopmentLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Simple logging
	logger.Info("Application started")

	// Structured logging with fields
	logger.Info("User logged in",
		golog.Field("user_id", 12345),
		golog.Field("username", "john_doe"),
		golog.Field("ip", "192.168.1.1"),
	)

	// Different log levels
	logger.Debug("Debug information", golog.Field("debug_key", "debug_value"))
	logger.Warn("Warning message", golog.Field("warning_code", 1001))
	logger.Error("Error occurred", golog.Field("error_code", 5001))

	// Child logger with common fields
	requestLogger := logger.With(
		golog.Field("request_id", "abc-123"),
		golog.Field("service", "api"),
	)

	requestLogger.Info("Processing request")
	requestLogger.Info("Request completed", golog.Field("duration_ms", 45))
}
