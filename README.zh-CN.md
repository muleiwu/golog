# golog

[![Go Reference](https://pkg.go.dev/badge/github.com/muleiwu/golog.svg)](https://pkg.go.dev/github.com/muleiwu/golog)
[![Go Report Card](https://goreportcard.com/badge/github.com/muleiwu/golog)](https://goreportcard.com/report/github.com/muleiwu/golog)

[English](README.md) | [中文](README.zh-CN.md)

一个灵活且结构化的 Go 日志库,基于 [uber-go/zap](https://github.com/uber-go/zap) 构建,实现了 [gsr](https://github.com/muleiwu/gsr) 日志接口。

## 特性

- 🚀 **高性能**: 基于 uber-go/zap 构建,最快的结构化日志库之一
- 🎯 **结构化日志**: 支持强类型的结构化日志字段
- 🔧 **灵活配置**: 针对不同环境提供多种初始化选项
- 📊 **多日志级别**: Debug、Info、Notice、Warn、Error、Fatal 和 Panic
- 🎨 **多输出格式**: JSON 和控制台编码
- 🔌 **接口兼容**: 实现 gsr.Logger 接口
- 🛠️ **易于使用**: 简单直观的 API
- 📍 **准确的调用信息**: 日志显示正确的调用文件和行号位置

## 安装

```bash
go get github.com/muleiwu/golog
```

## 快速开始

```go
package main

import (
    "github.com/muleiwu/golog"
)

func main() {
    // 创建开发环境日志器
    logger, err := golog.NewDevelopmentLogger()
    if err != nil {
        panic(err)
    }
    defer logger.Sync()

    // 简单日志记录
    logger.Info("应用程序已启动")

    // 带字段的结构化日志
    logger.Info("用户登录",
        golog.Field("user_id", 12345),
        golog.Field("username", "john_doe"),
        golog.Field("ip", "192.168.1.1"),
    )
}
```

## 使用方法

### 日志器初始化

#### 开发环境日志器

最适合开发环境,具有人类可读的控制台输出:

```go
logger, err := golog.NewDevelopmentLogger()
if err != nil {
    panic(err)
}
defer logger.Sync()
```

#### 生产环境日志器

针对生产环境优化,使用 JSON 输出:

```go
logger, err := golog.NewProductionLogger()
if err != nil {
    panic(err)
}
defer logger.Sync()
```

#### 示例日志器

仅用于测试目的(不推荐用于生产环境):

```go
logger := golog.NewLogger()
```

#### 自定义配置

使用自定义设置创建日志器:

```go
logger, err := golog.NewLoggerWithConfig(golog.Config{
    Level:            golog.DebugLevel,  // 使用 golog.Level 常量
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

**可用的日志级别:**
- `golog.DebugLevel` - 调试消息
- `golog.InfoLevel` - 信息消息（默认）
- `golog.WarnLevel` - 警告消息
- `golog.ErrorLevel` - 错误消息
- `golog.FatalLevel` - 致命消息（调用 os.Exit）
- `golog.PanicLevel` - 恐慌消息（记录后 panic）

#### 从已有的 Zap 日志器创建

包装现有的 zap.Logger:

```go
zapLogger, _ := zap.NewProduction()
logger := golog.NewLoggerWithZap(zapLogger)
```

### 日志级别

```go
logger.Debug("调试消息", golog.Field("key", "value"))
logger.Info("信息消息", golog.Field("key", "value"))
logger.Notice("通知消息", golog.Field("key", "value"))  // 映射到 Info
logger.Warn("警告消息", golog.Field("key", "value"))
logger.Error("错误消息", golog.Field("key", "value"))
logger.Fatal("致命错误消息", golog.Field("key", "value"))    // 调用 os.Exit(1)
logger.Panic("恐慌消息", golog.Field("key", "value"))      // 记录后触发 panic
```

### 结构化日志

通过字段为日志添加上下文:

```go
logger.Info("处理请求",
    golog.Field("request_id", "abc-123"),
    golog.Field("method", "GET"),
    golog.Field("path", "/api/users"),
    golog.Field("duration_ms", 45),
)
```

### 子日志器

创建带有预填充字段的子日志器:

```go
// 创建带有公共字段的子日志器
requestLogger := logger.With(
    golog.Field("request_id", "abc-123"),
    golog.Field("user_id", 12345),
)

// requestLogger 的所有日志都将包含这些字段
requestLogger.Info("请求开始")
requestLogger.Info("请求完成")
```

### 高级用法

#### 直接访问 Zap

访问底层的 zap.Logger 以使用高级特性:

```go
zapLogger := logger.GetZapLogger()
// 使用 zap 特定功能
```

#### 直接使用 Zap 字段

为了更好的性能,可以直接使用 zap 字段:

```go
import "go.uber.org/zap"

childLogger := logger.WithZapFields(
    zap.String("service", "api"),
    zap.Int("port", 8080),
)
```

## 配置选项

`Config` 结构体支持以下选项:

| 字段 | 类型 | 描述 |
|------|------|------|
| `Level` | `golog.Level` | 最小日志级别 (DebugLevel、InfoLevel、WarnLevel、ErrorLevel、FatalLevel、PanicLevel) |
| `Development` | `bool` | 启用开发模式(更易读) |
| `Encoding` | `string` | 输出格式: "json" 或 "console" |
| `OutputPaths` | `[]string` | 输出目标(如 "stdout"、文件路径) |
| `ErrorOutputPaths` | `[]string` | 错误输出目标(如 "stderr") |
| `CallerSkip` | `uint` | 额外跳过的栈帧数(自动 +1 用于 golog)。默认值：0(总共跳过 1 层)。单层封装设为 1，双层封装设为 2，以此类推。 |

### 日志级别说明

- `DebugLevel`: 细粒度的调试信息
- `InfoLevel`: 一般信息消息
- `WarnLevel`: 潜在有害情况的警告消息
- `ErrorLevel`: 严重问题的错误消息
- `FatalLevel`: 导致程序退出的严重错误
- `PanicLevel`: 导致 panic 的严重错误

## 最佳实践

1. **始终调用 `Sync()`**: 确保程序退出前日志被刷新
   ```go
   defer logger.Sync()  // 安全调用，自动处理 stdout/stderr
   ```

   注意：`Sync()` 会自动忽略 stdout/stderr 的错误（在某些系统上无法同步），所以你可以安全地使用 `defer logger.Sync()` 而不用担心 "bad file descriptor" 错误。

2. **使用适当的日志级别**:
   - `Debug` 用于开发调试
   - `Info` 用于一般信息
   - `Warn` 用于潜在有害情况
   - `Error` 用于需要关注的错误
   - `Fatal`/`Panic` 仅用于关键故障

3. **使用结构化字段**: 不要使用字符串格式化,而是使用字段
   ```go
   // 推荐
   logger.Info("用户操作", golog.Field("user_id", userID), golog.Field("action", "login"))

   // 避免
   logger.Info(fmt.Sprintf("用户 %d 执行操作: login", userID))
   ```

4. **创建子日志器**: 用于请求范围或特定上下文的日志记录
   ```go
   requestLogger := logger.With(golog.Field("request_id", requestID))
   ```

5. **在生产环境使用生产日志器**: 开发日志器未针对性能优化

## 示例

### Web 服务器示例

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
        logger.Info("收到请求",
            golog.Field("method", r.Method),
            golog.Field("path", r.URL.Path),
            golog.Field("remote_addr", r.RemoteAddr),
        )
        w.Write([]byte("Hello, World!"))
    })

    logger.Info("服务器启动", golog.Field("port", 8080))
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Fatal("服务器启动失败", golog.Field("error", err))
    }
}
```

### 错误处理示例

```go
func processUser(logger *golog.Logger, userID int) error {
    userLogger := logger.With(golog.Field("user_id", userID))

    userLogger.Debug("开始处理用户")

    user, err := fetchUser(userID)
    if err != nil {
        userLogger.Error("获取用户失败", golog.Field("error", err))
        return err
    }

    userLogger.Info("成功获取用户", golog.Field("username", user.Name))
    return nil
}
```

### 封装日志器示例

如果你在自己的日志器中封装 golog，将 `CallerSkip` 设置为封装层数：

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
        CallerSkip:  1, // 1 层封装(自动 +1 给 golog = 总共 2 层)
    })
    if err != nil {
        return nil, err
    }
    return &MyLogger{logger: logger}, nil
}

func (m *MyLogger) Info(msg string, fields ...gsr.LoggerField) {
    m.logger.Info(msg, fields...)  // 现在显示正确的调用位置
}
```

**CallerSkip 设置指南:**
- `0` = 直接使用(跳过 1 层: 仅 golog) - 与预设日志器相同
- `1` = 单层封装(跳过 2 层: golog + 你的封装)
- `2` = 双层封装(跳过 3 层: golog + 2 层封装)
- `N` = N 层封装(跳过 N+1 层: golog + N 层封装)

## 依赖

- [go.uber.org/zap](https://github.com/uber-go/zap) - 快速的结构化日志库
- [github.com/muleiwu/gsr](https://github.com/muleiwu/gsr) - 日志接口定义

## 贡献

欢迎贡献!请随时提交 Pull Request。

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

- 使用 [uber-go/zap](https://github.com/uber-go/zap) 构建
- 实现 [gsr](https://github.com/muleiwu/gsr) 日志接口