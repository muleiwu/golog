package main

import (
	"github.com/muleiwu/golog"
)

func main() {
	// 默认配置：显示短路径
	logger1, err := golog.NewLoggerWithConfig(golog.Config{
		Level:             golog.InfoLevel,
		Development:       true,
		Encoding:          "console",
		OutputPaths:       []string{"stdout"},
		DisableCallerTrim: false, // 默认：短路径
	})
	if err != nil {
		panic(err)
	}
	defer logger1.Sync()

	logger1.Info("默认显示短路径")
	logger1.Info("例如: main.go:25", golog.Field("path_type", "short"))

	println("\n--- 分隔线 ---\n")

	// 完整路径配置：显示从模块根目录开始的完整路径
	logger2, err := golog.NewLoggerWithConfig(golog.Config{
		Level:             golog.InfoLevel,
		Development:       true,
		Encoding:          "console",
		OutputPaths:       []string{"stdout"},
		DisableCallerTrim: true, // 显示完整路径
	})
	if err != nil {
		panic(err)
	}
	defer logger2.Sync()

	logger2.Info("显示从模块根目录开始的完整路径")
	logger2.Info("例如: examples/full_caller_path/main.go:45", golog.Field("path_type", "full"))

	println("\n--- 说明 ---")
	println("DisableCallerTrim: false -> 适合开发环境，路径简短")
	println("DisableCallerTrim: true  -> 适合复杂项目，显示完整路径结构")
}
