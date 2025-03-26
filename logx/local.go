package logx

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewFileWriter 创建一个文件写入器，支持日志轮转
func NewFileWriter(filename string) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,   // 每个日志文件最大 10MB
		MaxBackups: 3,    // 保留最多 3 个备份
		MaxAge:     28,   // 保留日志 28 天
		Compress:   true, // 是否压缩
	}
	return zapcore.AddSync(lumberjackLogger)
}
