package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFile = "server.log"
)

var logger *zap.Logger

// Initialize the global logger
func InitLogger(isDebug bool) *zap.Logger {
	var cfg zap.Config

	if isDebug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.OutputPaths = []string{
		logFile,
		"stdout",
	}

	var err error
	logger, err = cfg.Build()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		panic(err)
	}

	logger.Info("Logger initialized", zap.Bool("isDebug", isDebug))
	return logger
}

func Sync() {
	if err := logger.Sync(); err != nil {
		fmt.Printf("Failed to ensure all buffered log etries are flushed to the output: %s\n", err)
	}
}

// Shorthand functions for different log levels
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
