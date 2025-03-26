package example

import (
	"fmt"
	"github.com/trancecho/open-sdk/logx"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志记录器
	err := logx.InitLogger(
		"amqp://root:123456@localhost:56720/",
		"logs",
		"app.log",
	)
	if err != nil {
		// 如果初始化失败，使用标准库记录错误并退出
		panic(err)
	}
	// 确保在程序退出时同步日志
	defer logx.SyncLogger()

	// 示例登录和注册日志
	username := "john_doe"

	// 登录日志
	logx.Info("User login",
		zap.String("username", username),
		zap.String("ip_address", "192.168.1.100"),
		zap.String("status", "success"),
	)

	// 注册日志
	logx.Info("User registration",
		zap.String("username", username),
		zap.String("email", "john@example.com"),
		zap.String("referral_code", "REF123"),
	)

	// 其他日志示例
	logx.Debug("This is a debug message", zap.Int("debug_level", 1))
	logx.Warn("This is a warning message", zap.String("module", "auth"))
	logx.Error("This is an error message", zap.Error(fmt.Errorf("sample error")))
	// logx.Fatal("This is a fatal message") // 会退出程序
	// logx.Panic("This is a panic message") // 会抛出异常
}
