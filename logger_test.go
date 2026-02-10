package golog

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	if logger.logger == nil {
		t.Fatal("Internal zap logger is nil")
	}
}

func TestNewDevelopmentLogger(t *testing.T) {
	logger, err := NewDevelopmentLogger()
	if err != nil {
		t.Fatalf("NewDevelopmentLogger failed: %v", err)
	}
	if logger == nil {
		t.Fatal("NewDevelopmentLogger returned nil")
	}
	if logger.logger == nil {
		t.Fatal("Internal zap logger is nil")
	}
}

func TestNewProductionLogger(t *testing.T) {
	logger, err := NewProductionLogger()
	if err != nil {
		t.Fatalf("NewProductionLogger failed: %v", err)
	}
	if logger == nil {
		t.Fatal("NewProductionLogger returned nil")
	}
	if logger.logger == nil {
		t.Fatal("Internal zap logger is nil")
	}
	defer logger.Sync()
}

func TestNewLoggerWithConfig(t *testing.T) {
	config := Config{
		Level:            zapcore.DebugLevel,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := NewLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}
	if logger == nil {
		t.Fatal("NewLoggerWithConfig returned nil")
	}
	if logger.logger == nil {
		t.Fatal("Internal zap logger is nil")
	}
	defer logger.Sync()
}

func TestLoggerLevels(t *testing.T) {
	logger := NewLogger()
	defer logger.Sync()

	// Test all log levels (except Fatal and Panic which would terminate/panic)
	logger.Debug("debug message", Field("key", "value"))
	logger.Info("info message", Field("key", "value"))
	logger.Notice("notice message", Field("key", "value"))
	logger.Warn("warn message", Field("key", "value"))
	logger.Error("error message", Field("key", "value"))

	// Test without fields
	logger.Info("message without fields")
}

func TestLoggerWith(t *testing.T) {
	logger := NewLogger()
	defer logger.Sync()

	childLogger := logger.With(
		Field("service", "test"),
		Field("version", "1.0"),
	)

	if childLogger == nil {
		t.Fatal("With returned nil")
	}
	if childLogger.logger == nil {
		t.Fatal("Child logger's internal zap logger is nil")
	}

	// Log with child logger
	childLogger.Info("message from child logger", Field("extra", "field"))
}

func TestLoggerGetZapLogger(t *testing.T) {
	logger := NewLogger()
	zapLogger := logger.GetZapLogger()

	if zapLogger == nil {
		t.Fatal("GetZapLogger returned nil")
	}
}

func TestField(t *testing.T) {
	field := Field("test_key", "test_value")

	if field == nil {
		t.Fatal("Field returned nil")
	}

	loggerField, ok := field.(*LoggerField)
	if !ok {
		t.Fatal("Field did not return *LoggerField")
	}

	if loggerField.GetKey() != "test_key" {
		t.Errorf("Expected key 'test_key', got '%s'", loggerField.GetKey())
	}

	if loggerField.GetValue() != "test_value" {
		t.Errorf("Expected value 'test_value', got '%v'", loggerField.GetValue())
	}
}

func TestMultipleFields(t *testing.T) {
	logger := NewLogger()
	defer logger.Sync()

	logger.Info("testing multiple fields",
		Field("string", "value"),
		Field("int", 42),
		Field("bool", true),
		Field("float", 3.14),
	)
}

func BenchmarkLoggerInfo(b *testing.B) {
	logger := NewLogger()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", Field("iteration", i))
	}
}

func BenchmarkLoggerInfoWithMultipleFields(b *testing.B) {
	logger := NewLogger()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message",
			Field("iteration", i),
			Field("name", "test"),
			Field("active", true),
			Field("score", 95.5),
		)
	}
}
