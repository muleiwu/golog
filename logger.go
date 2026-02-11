package golog

import (
	"github.com/muleiwu/gsr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents the logging level
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
)

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	default:
		return "unknown"
	}
}

// toZapLevel converts our Level to zapcore.Level
func (l Level) toZapLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// Logger wraps zap.Logger and implements the gsr.Logger interface
type Logger struct {
	logger *zap.Logger
}

// Config holds the configuration for creating a new logger
type Config struct {
	// Level sets the minimum enabled logging level
	Level Level
	// Development puts the logger in development mode
	Development bool
	// Encoding sets the logger's encoding (json or console)
	Encoding string
	// OutputPaths is a list of URLs or file paths to write logging output to
	OutputPaths []string
	// ErrorOutputPaths is a list of URLs to write internal logger errors to
	ErrorOutputPaths []string
	// CallerSkip increases the number of callers skipped by caller annotation
	// Default is 0, which will be automatically set to 1 (skip golog wrapper).
	// Set to 1+ if you wrap golog in your own logger (1 = single wrap, 2 = double wrap, etc.).
	CallerSkip uint
}

// NewLogger creates a new logger with example configuration (for testing only)
func NewLogger() *Logger {
	// Add caller skip to show correct file and line number
	logger := zap.NewExample(zap.AddCallerSkip(1))
	return &Logger{
		logger: logger,
	}
}

// NewDevelopmentLogger creates a logger suitable for development
// with human-readable console output and debug-level logging
func NewDevelopmentLogger() (*Logger, error) {
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger}, nil
}

// NewProductionLogger creates a logger suitable for production
// with JSON output and info-level logging
func NewProductionLogger() (*Logger, error) {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger}, nil
}

// NewLoggerWithConfig creates a new logger with custom configuration
func NewLoggerWithConfig(config Config) (*Logger, error) {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(config.Level.toZapLevel()),
		Development:      config.Development,
		Encoding:         config.Encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      config.OutputPaths,
		ErrorOutputPaths: config.ErrorOutputPaths,
	}

	if config.Development {
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	// Use configured caller skip + 1 (for golog wrapper)
	// If CallerSkip is 0 (default), this results in 1 (same as preset loggers)
	// If CallerSkip is 1+, it adds 1 for the golog wrapper layer
	callerSkip := config.CallerSkip

	logger, err := zapConfig.Build(zap.AddCallerSkip(int(callerSkip + 1)))
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger}, nil
}

// NewLoggerWithZap creates a logger from an existing zap.Logger
// Note: If you need caller skip, pass a logger with AddCallerSkip already configured
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
