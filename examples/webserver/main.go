package main

import (
	"net/http"
	"time"

	"github.com/muleiwu/golog"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Create a production-ready logger with custom configuration
	logger, err := golog.NewLoggerWithConfig(golog.Config{
		Level:            zapcore.InfoLevel,
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// HTTP handler with logging
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a request-scoped logger
		requestLogger := logger.With(
			golog.Field("method", r.Method),
			golog.Field("path", r.URL.Path),
			golog.Field("remote_addr", r.RemoteAddr),
		)

		requestLogger.Info("Request received")

		// Simulate some processing
		w.Write([]byte("Hello, World!"))

		// Log completion with duration
		requestLogger.Info("Request completed",
			golog.Field("duration_ms", time.Since(start).Milliseconds()),
			golog.Field("status", http.StatusOK),
		)
	})

	logger.Info("Starting HTTP server", golog.Field("port", 8080))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Failed to start server", golog.Field("error", err))
	}
}
