package golog

import (
	"github.com/muleiwu/gsr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger and implements the gsr.Logger interface
type Logger struct {
	logger *zap.Logger
}

// Config holds the configuration for creating a new logger
type Config struct {
	// Level sets the minimum enabled logging level
	Level zapcore.Level
	// Development puts the logger in development mode
	Development bool
	// Encoding sets the logger's encoding (json or console)
	Encoding string
	// OutputPaths is a list of URLs or file paths to write logging output to
	OutputPaths []string
	// ErrorOutputPaths is a list of URLs to write internal logger errors to
	ErrorOutputPaths []string
}

// NewLogger creates a new logger with example configuration (for testing only)
func NewLogger() *Logger {
	return &Logger{
		logger: zap.NewExample(),
	}
}

// NewDevelopmentLogger creates a logger suitable for development
// with human-readable console output and debug-level logging
func NewDevelopmentLogger() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger}, nil
}

// NewProductionLogger creates a logger suitable for production
// with JSON output and info-level logging
func NewProductionLogger() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger}, nil
}

// NewLoggerWithConfig creates a new logger with custom configuration
func NewLoggerWithConfig(config Config) (*Logger, error) {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(config.Level),
		Development:      config.Development,
		Encoding:         config.Encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      config.OutputPaths,
		ErrorOutputPaths: config.ErrorOutputPaths,
	}

	if config.Development {
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger}, nil
}

// NewLoggerWithZap creates a logger from an existing zap.Logger
func NewLoggerWithZap(zapLogger *zap.Logger) *Logger {
	return &Logger{logger: zapLogger}
}

// getFields converts gsr.LoggerField to zap.Field
func (l *Logger) getFields(args ...gsr.LoggerField) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		fields = append(fields, zap.Any(arg.GetKey(), arg.GetValue()))
	}

	return fields
}

// Debug logs a message at DebugLevel
func (l *Logger) Debug(format string, args ...gsr.LoggerField) {
	l.logger.Debug(format, l.getFields(args...)...)
}

// Info logs a message at InfoLevel
func (l *Logger) Info(format string, args ...gsr.LoggerField) {
	l.logger.Info(format, l.getFields(args...)...)
}

// Notice logs a message at InfoLevel (alias for Info)
// Notice level is mapped to Info as zap doesn't have a separate Notice level
func (l *Logger) Notice(format string, args ...gsr.LoggerField) {
	l.logger.Info(format, l.getFields(args...)...)
}

// Warn logs a message at WarnLevel
func (l *Logger) Warn(format string, args ...gsr.LoggerField) {
	l.logger.Warn(format, l.getFields(args...)...)
}

// Error logs a message at ErrorLevel
func (l *Logger) Error(format string, args ...gsr.LoggerField) {
	l.logger.Error(format, l.getFields(args...)...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func (l *Logger) Fatal(format string, args ...gsr.LoggerField) {
	l.logger.Fatal(format, l.getFields(args...)...)
}

// Panic logs a message at PanicLevel and then panics
func (l *Logger) Panic(format string, args ...gsr.LoggerField) {
	l.logger.Panic(format, l.getFields(args...)...)
}

// Sync flushes any buffered log entries
// Applications should call Sync before exiting
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// With creates a child logger with additional fields
func (l *Logger) With(args ...gsr.LoggerField) *Logger {
	return &Logger{
		logger: l.logger.With(l.getFields(args...)...),
	}
}

// WithZapFields creates a child logger with additional zap fields
func (l *Logger) WithZapFields(fields ...zap.Field) *Logger {
	return &Logger{
		logger: l.logger.With(fields...),
	}
}

// GetZapLogger returns the underlying zap.Logger
// This is useful when you need direct access to zap features
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.logger
}
