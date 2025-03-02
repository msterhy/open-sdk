package logx

import (
	"github.com/trancecho/open-sdk/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// NameSpace - 提供带有模块命名空间的logger
func NameSpace(name string) *zap.SugaredLogger {
	return zap.S().Named(name)
}

func getLogWriter() zapcore.WriteSyncer {
	if config.GetConfig().LogPath == "" {
		config.GetConfig().LogPath = "app.log"
		print("LogPath 未设置, 使用默认值app.log\n")
	}
	lj := &lumberjack.Logger{
		Filename:   config.GetConfig().LogPath,
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lj)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func InitLogger(level zapcore.LevelEnabler) {
	writeSyncer := getLogWriter()
	if level == zapcore.DebugLevel {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)))
}

// Debug logs a message at debug level.
func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

// Warn logs a message at warn level.
func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

// DPanic logs a message at DPanic level. DPanic level logs are particularly
// important errors. In development the logger panics after logging them.
func DPanic(msg string, fields ...zap.Field) {
	zap.L().DPanic(msg, fields...)
}

// Panic logs a message at panic level. The logger then panics.
func Panic(msg string, fields ...zap.Field) {
	zap.L().Panic(msg, fields...)
}

// Fatal logs a message at fatal level. The logger then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}
