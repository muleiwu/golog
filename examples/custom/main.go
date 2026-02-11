package main

import (
	"github.com/muleiwu/golog"
)

func main() {
	// Custom logger configuration
	logger, err := golog.NewLoggerWithConfig(golog.Config{
		Level:       golog.DebugLevel,
		Development: true,
		Encoding:    "console",
		OutputPaths: []string{
			"stdout",
			"/tmp/app.log", // Also write to file
		},
		ErrorOutputPaths: []string{"stderr"},
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Custom logger initialized")

	// Demonstrate all log levels
	logger.Debug("This is a debug message",
		golog.Field("environment", "development"),
	)

	logger.Info("Application is running",
		golog.Field("version", "1.0.0"),
		golog.Field("build", "2024-02-10"),
	)

	logger.Warn("Deprecated feature used",
		golog.Field("feature", "old_api"),
		golog.Field("replacement", "new_api_v2"),
	)

	logger.Error("Failed to process item",
		golog.Field("item_id", 999),
		golog.Field("reason", "invalid format"),
	)

	// Create child logger with service context
	serviceLogger := logger.With(
		golog.Field("service", "payment"),
		golog.Field("version", "2.0"),
	)

	serviceLogger.Info("Payment service started")
	serviceLogger.Info("Processing payment",
		golog.Field("amount", 99.99),
		golog.Field("currency", "USD"),
	)
}
