package utils

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel string = "DEBUG"
	InfoLevel  string = "INFO"
	WarnLevel  string = "WARN"
	ErrorLevel string = "ERROR"
	FatalLevel string = "FATAL"
)

var (
	loggerInstance *zap.Logger
)

// Initialize logger at package level, so it runs once when the package is imported
func init() {
	var err error
	loggerInstance, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	loggerInstance.Info("Logger is initialized")
}

func Logger(level string, msg string, fields ...zapcore.Field) {
	if loggerInstance == nil {
		log.Fatalf("Logger is not initialized")
	}
	switch level {
	case DebugLevel:
		loggerInstance.Debug(msg, fields...)
	case WarnLevel:
		loggerInstance.Warn(msg, fields...)
	case ErrorLevel:
		loggerInstance.Error(msg, fields...)
	default:
		loggerInstance.Info(msg, fields...)
	}
}
