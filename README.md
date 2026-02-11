# golog

[![Go Reference](https://pkg.go.dev/badge/github.com/muleiwu/golog.svg)](https://pkg.go.dev/github.com/muleiwu/golog)
[![Go Report Card](https://goreportcard.com/badge/github.com/muleiwu/golog)](https://goreportcard.com/report/github.com/muleiwu/golog)

[English](README.md) | [中文](README.zh-CN.md)

A flexible and structured logging library for Go, built on top of [uber-go/zap](https://github.com/uber-go/zap) and implementing the [gsr](https://github.com/muleiwu/gsr) logger interface.

## Features

- 🚀 **High Performance**: Built on uber-go/zap, one of the fastest structured logging libraries
- 🎯 **Structured Logging**: Support for strongly-typed, structured log fields
- 🔧 **Flexible Configuration**: Multiple initialization options for different environments
- 📊 **Multiple Log Levels**: Debug, Info, Notice, Warn, Error, Fatal, and Panic
- 🎨 **Multiple Output Formats**: JSON and console encoding
- 🔌 **Interface Compliant**: Implements the gsr.Logger interface
- 🛠️ **Easy to Use**: Simple and intuitive API
- 📍 **Accurate Caller Info**: Logs show the correct file and line number where the log was called

## Installation

```bash
go get github.com/muleiwu/golog
```

## Quick Start

```go
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
}
```

## Usage

### Logger Initialization

#### Development Logger

Best for development environments with human-readable console output:

```go
logger, err := golog.NewDevelopmentLogger()
if err != nil {
    panic(err)
}
defer logger.Sync()
```

#### Production Logger

Optimized for production with JSON output:

```go
logger, err := golog.NewProductionLogger()
if err != nil {
    panic(err)
}
defer logger.Sync()
```

#### Example Logger

For testing purposes only (not recommended for production):

```go
logger := golog.NewLogger()
```

#### Custom Configuration

Create a logger with custom settings:

```go
logger, err := golog.NewLoggerWithConfig(golog.Config{
    Level:            golog.DebugLevel,  // Use golog.Level constants
    Development:      true,
    Encoding:         "console",
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    ErrorOutputPaths: []string{"stderr"},
})
if err != nil {
    panic(err)
}
defer logger.Sync()
```

**Available Log Levels:**
- `golog.DebugLevel` - Debug messages
- `golog.InfoLevel` - Informational messages (default)
- `golog.WarnLevel` - Warning messages
- `golog.ErrorLevel` - Error messages
- `golog.FatalLevel` - Fatal messages (calls os.Exit)
- `golog.PanicLevel` - Panic messages (panics after logging)

#### From Existing Zap Logger

Wrap an existing zap.Logger:

```go
zapLogger, _ := zap.NewProduction()
logger := golog.NewLoggerWithZap(zapLogger)
```

### Logging Levels

```go
logger.Debug("Debug message", golog.Field("key", "value"))
logger.Info("Info message", golog.Field("key", "value"))
logger.Notice("Notice message", golog.Field("key", "value"))  // Mapped to Info
logger.Warn("Warning message", golog.Field("key", "value"))
logger.Error("Error message", golog.Field("key", "value"))
logger.Fatal("Fatal message", golog.Field("key", "value"))    // Calls os.Exit(1)
logger.Panic("Panic message", golog.Field("key", "value"))    // Panics after logging
```

### Structured Logging

Add context to your logs with fields:

```go
logger.Info("Processing request",
    golog.Field("request_id", "abc-123"),
    golog.Field("method", "GET"),
    golog.Field("path", "/api/users"),
    golog.Field("duration_ms", 45),
)
```

### Child Loggers

Create child loggers with pre-populated fields:

```go
// Create a child logger with common fields
requestLogger := logger.With(
    golog.Field("request_id", "abc-123"),
    golog.Field("user_id", 12345),
)

// All logs from requestLogger will include these fields
requestLogger.Info("Request started")
requestLogger.Info("Request completed")
```

### Advanced Usage

#### Direct Zap Access

Access the underlying zap.Logger for advanced features:

```go
zapLogger := logger.GetZapLogger()
// Use zap-specific features
```

#### Using Zap Fields Directly

For better performance, you can use zap fields directly:

```go
import "go.uber.org/zap"

childLogger := logger.WithZapFields(
    zap.String("service", "api"),
    zap.Int("port", 8080),
)
```

## Configuration Options

The `Config` struct supports the following options:

| Field | Type | Description |
|-------|------|-------------|
| `Level` | `golog.Level` | Minimum logging level (DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel) |
| `Development` | `bool` | Enable development mode (more human-readable) |
| `Encoding` | `string` | Output format: "json" or "console" |
| `OutputPaths` | `[]string` | Output destinations (e.g., "stdout", file paths) |
| `ErrorOutputPaths` | `[]string` | Error output destinations (e.g., "stderr") |
| `CallerSkip` | `uint` | Additional stack frames to skip (automatically +1 for golog). Default: 0 (total skip=1). Set to 1 for single wrapper, 2 for double wrapper, etc. |
| `DisableCallerTrim` | `bool` | Disable trimming of caller path. Default: false (shows short path like `service/server.go:67`). Set to true for full path from module root (like `pkg/service/cron/service/server.go:67`) |

### Log Levels

- `DebugLevel`: Fine-grained debugging information
- `InfoLevel`: General informational messages
- `WarnLevel`: Warning messages for potentially harmful situations
- `ErrorLevel`: Error messages for serious problems
- `FatalLevel`: Very severe errors that will lead to program exit
- `PanicLevel`: Very severe errors that will cause a panic

## Best Practices

1. **Always call `Sync()`**: Ensure logs are flushed before program exit
   ```go
   defer logger.Sync()  // Safe to call, handles stdout/stderr gracefully
   ```

   Note: `Sync()` automatically ignores errors from stdout/stderr (which cannot be synced on some systems), so you can safely use `defer logger.Sync()` without worrying about "bad file descriptor" errors.

2. **Use appropriate log levels**:
   - `Debug` for development debugging
   - `Info` for general information
   - `Warn` for potentially harmful situations
   - `Error` for errors that need attention
   - `Fatal`/`Panic` only for critical failures

3. **Use structured fields**: Instead of string formatting, use fields
   ```go
   // Good
   logger.Info("User action", golog.Field("user_id", userID), golog.Field("action", "login"))

   // Avoid
   logger.Info(fmt.Sprintf("User %d performed action: login", userID))
   ```

4. **Create child loggers**: For request-scoped or context-specific logging
   ```go
   requestLogger := logger.With(golog.Field("request_id", requestID))
   ```

5. **Use production logger in production**: Development logger is not optimized for performance

## Examples

### Web Server Example

```go
package main

import (
    "net/http"
    "github.com/muleiwu/golog"
    "go.uber.org/zap/zapcore"
)

func main() {
    logger, err := golog.NewLoggerWithConfig(golog.Config{
        Level:            golog.InfoLevel,
        Development:      false,
        Encoding:         "json",
        OutputPaths:      []string{"stdout", "/var/log/server.log"},
        ErrorOutputPaths: []string{"stderr"},
    })
    if err != nil {
        panic(err)
    }
    defer logger.Sync()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logger.Info("Request received",
            golog.Field("method", r.Method),
            golog.Field("path", r.URL.Path),
            golog.Field("remote_addr", r.RemoteAddr),
        )
        w.Write([]byte("Hello, World!"))
    })

    logger.Info("Server starting", golog.Field("port", 8080))
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Fatal("Server failed to start", golog.Field("error", err))
    }
}
```

### Error Handling Example

```go
func processUser(logger *golog.Logger, userID int) error {
    userLogger := logger.With(golog.Field("user_id", userID))

    userLogger.Debug("Starting user processing")

    user, err := fetchUser(userID)
    if err != nil {
        userLogger.Error("Failed to fetch user", golog.Field("error", err))
        return err
    }

    userLogger.Info("User fetched successfully", golog.Field("username", user.Name))
    return nil
}
```

### Wrapping Logger Example

If you wrap golog in your own logger, set `CallerSkip` to the number of wrapper layers:

```go
type MyLogger struct {
    logger *golog.Logger
}

func NewMyLogger() (*MyLogger, error) {
    logger, err := golog.NewLoggerWithConfig(golog.Config{
        Level:       golog.InfoLevel,
        Development: true,
        Encoding:    "console",
        OutputPaths: []string{"stdout"},
        CallerSkip:  1, // 1 wrapper layer (auto +1 for golog = total 2)
    })
    if err != nil {
        return nil, err
    }
    return &MyLogger{logger: logger}, nil
}

func (m *MyLogger) Info(msg string, fields ...gsr.LoggerField) {
    m.logger.Info(msg, fields...)  // Now shows correct caller location
}
```

**CallerSkip Guide:**
- `0` = Direct usage (skip 1: golog only) - same as preset loggers
- `1` = Single wrapper (skip 2: golog + your wrapper)
- `2` = Double wrapper (skip 3: golog + 2 wrappers)
- `N` = N wrappers (skip N+1: golog + N wrappers)

### Full Caller Path Example

For complex projects with deep directory structures, use `DisableCallerTrim`:

```go
logger, err := golog.NewLoggerWithConfig(golog.Config{
    Level:             golog.InfoLevel,
    Development:       true,
    Encoding:          "console",
    OutputPaths:       []string{"stdout"},
    DisableCallerTrim: true,  // Show full path from module root
})

// Output: 2026-02-11T10:09:15.151+0800    INFO    pkg/service/cron/service/cron_server.go:67    正在停止服务...
// Instead of: service/cron_server.go:67
```

**When to use:**
- ✅ Complex projects with deep directory structures
- ✅ Multiple packages with similar file names
- ✅ Need full context of where the log came from
- ❌ Simple projects (default short path is cleaner)

## Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap) - Fast, structured logging
- [github.com/muleiwu/gsr](https://github.com/muleiwu/gsr) - Logger interface definition

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [uber-go/zap](https://github.com/uber-go/zap)
- Implements [gsr](https://github.com/muleiwu/gsr) logger interface