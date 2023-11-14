package logger

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfg    config.AppConfiguration
	zapLog *zap.Logger
)

func init() {
	var err error
	if cfg.Mode == "production" {
		config := zap.NewProductionConfig()
		enccoderConfig := zap.NewProductionEncoderConfig()
		zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
		enccoderConfig.StacktraceKey = "" // to hide stacktrace info
		config.EncoderConfig = enccoderConfig

		zapLog, err = config.Build(zap.AddCallerSkip(1))
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLog, err = config.Build()
	}

	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Panic(message string, fields ...zap.Field) {
	zapLog.Panic(message, fields...)
}

func DPanic(message string, fields ...zap.Field) {
	zapLog.DPanic(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}
