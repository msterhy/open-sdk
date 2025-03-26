package logx

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

// logger 是私有的 Zap 日志记录器
var logger *zap.Logger
var once sync.Once
var initErr error

// InitLogger 初始化全局日志记录器
func InitLogger(amqpURL, queueName, filename string) error {
	once.Do(func() {
		// 配置 RabbitMQ
		rabbitWriter, err := NewRabbitMQWriter(amqpURL, queueName)
		if err != nil {
			initErr = fmt.Errorf("failed to create RabbitMQ writer: %v", err)
			return
		}

		// 配置本地文件日志
		fileWriter := NewFileWriter(filename)

		// 使用 MultiWriteSyncer 同时写入 RabbitMQ 和本地文件
		multiWriter := zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(rabbitWriter),
			fileWriter,
		)

		// 配置 Zap 编码器
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 创建 Zap 核心
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			multiWriter,
			zap.InfoLevel,
		)

		// 创建 Logger
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

		// 记录初始化成功
		//logger.Info("Logger initialized",
		//	zap.String("amqp_url", amqpURL),
		//	zap.String("queue_name", queueName),
		//	zap.String("log_file", filename),
		//)
		log.Println("Logger initialized")
	})
	return initErr
}

// SyncLogger 同步日志（通常在应用程序退出时调用）
func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// 以下是简洁的日志接口

// Debug 记录 Debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Info 记录 Info 级别日志
func Info(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Warn 记录 Warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Error 记录 Error 级别日志
func Error(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

// Fatal 记录 Fatal 级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}

// Panic 记录 Panic 级别日志并抛出异常
func Panic(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Panic(msg, fields...)
	}
}
