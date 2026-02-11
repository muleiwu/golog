package main

import (
	"github.com/muleiwu/golog"
	"github.com/muleiwu/gsr"
)

// MyLogger 是对 golog.Logger 的二次封装
// MyLogger wraps golog.Logger for custom functionality
type MyLogger struct {
	logger *golog.Logger
}

// NewMyLogger 创建一个自定义的日志器
// 注意：因为又增加了一层包装，需要设置 CallerSkip=1
// CallerSkip 会自动 +1 来跳过 golog 层，所以 CallerSkip=1 表示额外跳过 1 层（MyLogger）
func NewMyLogger() (*MyLogger, error) {
	logger, err := golog.NewLoggerWithConfig(golog.Config{
		Level:            golog.InfoLevel,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		CallerSkip:       1, // 1 表示额外跳过 1 层封装（实际跳过 golog+MyLogger 共 2 层）
	})
	if err != nil {
		return nil, err
	}

	return &MyLogger{logger: logger}, nil
}

// Info 封装的 Info 方法
func (m *MyLogger) Info(msg string, fields ...gsr.LoggerField) {
	// 可以在这里添加自定义逻辑
	m.logger.Info(msg, fields...)
}

// Error 封装的 Error 方法
func (m *MyLogger) Error(msg string, fields ...gsr.LoggerField) {
	// 可以在这里添加自定义逻辑
	m.logger.Error(msg, fields...)
}

func (m *MyLogger) Sync() error {
	return m.logger.Sync()
}

func main() {
	// 使用封装的日志器
	myLogger, err := NewMyLogger()
	if err != nil {
		panic(err)
	}
	defer myLogger.Sync()

	// 这些日志应该显示 main.go 的行号，而不是 MyLogger 或 golog 的内部位置
	myLogger.Info("这是通过二次封装的日志器输出的")
	myLogger.Info("用户操作", golog.Field("user_id", 123), golog.Field("action", "login"))
	myLogger.Error("发生错误", golog.Field("error_code", 500))

	// 说明：
	// CallerSkip=0: 实际跳过 1 层（golog），会显示 wrapped_logger.go 的行号 ❌
	// CallerSkip=1: 实际跳过 2 层（golog+MyLogger），显示 main.go 的行号 ✅
	// CallerSkip=2: 实际跳过 3 层，会显示调用栈更上层的位置 ❌
}
