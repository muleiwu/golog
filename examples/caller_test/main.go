package main

import (
	"github.com/muleiwu/golog"
)

func main() {
	// Create a development logger to see caller info
	logger, err := golog.NewDevelopmentLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// These log calls should show THIS file and line numbers, not golog's internal files
	logger.Info("这条日志应该显示正确的文件名和行号")
	logger.Info("服务启动", golog.Field("port", 8080))
	logger.Warn("警告信息", golog.Field("code", 1001))
	logger.Error("错误信息", golog.Field("error", "something went wrong"))

	// Test with child logger
	serviceLogger := logger.With(
		golog.Field("service", "test"),
		golog.Field("version", "1.0"),
	)

	serviceLogger.Info("子日志器也应该显示正确的调用位置")
}
