package golog

import (
	"testing"
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
		Level:            DebugLevel,
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

func TestLevel(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{PanicLevel, "panic"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("Level(%d).String() = %v, want %v", tt.level, got, tt.expected)
		}
	}
}

func TestLevelInConfig(t *testing.T) {
	levels := []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}

	for _, level := range levels {
		config := Config{
			Level:            level,
			Development:      true,
			Encoding:         "console",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, err := NewLoggerWithConfig(config)
		if err != nil {
			t.Fatalf("NewLoggerWithConfig with Level=%v failed: %v", level, err)
		}
		defer logger.Sync()

		// Just verify it doesn't panic
		logger.Info("test message with level", Field("level", level.String()))
	}
}

func TestSync(t *testing.T) {
	// Test that Sync() doesn't return errors for stdout/stderr
	logger, err := NewLoggerWithConfig(Config{
		Level:            InfoLevel,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Info("Test message before sync")

	// This should not return an error even though stdout/stderr can't be synced
	err = logger.Sync()
	if err != nil {
		t.Errorf("Sync() returned unexpected error: %v", err)
	}
}

func TestSyncMultipleTimes(t *testing.T) {
	logger := NewLogger()

	logger.Info("Message 1")
	if err := logger.Sync(); err != nil {
		t.Errorf("First Sync() failed: %v", err)
	}

	logger.Info("Message 2")
	if err := logger.Sync(); err != nil {
		t.Errorf("Second Sync() failed: %v", err)
	}

	// Multiple syncs should be safe
	if err := logger.Sync(); err != nil {
		t.Errorf("Third Sync() failed: %v", err)
	}
}

func TestConfigWithCallerSkip(t *testing.T) {
	// Test default caller skip (0 means actual skip = 0+1 = 1, same as preset loggers)
	config1 := Config{
		Level:            DebugLevel,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		CallerSkip:       0, // 0 + 1 = 1 (default, same as NewDevelopmentLogger)
	}

	logger1, err := NewLoggerWithConfig(config1)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}
	defer logger1.Sync()

	// Test custom caller skip for wrapped loggers
	config2 := Config{
		Level:            InfoLevel,
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		CallerSkip:       1, // 1 + 1 = 2 (for single wrapper layer)
	}

	logger2, err := NewLoggerWithConfig(config2)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig with CallerSkip=1 failed: %v", err)
	}
	defer logger2.Sync()

	// Just verify they don't panic
	logger1.Info("test with default caller skip (0+1=1)")
	logger2.Info("test with custom caller skip (1+1=2)")
}
