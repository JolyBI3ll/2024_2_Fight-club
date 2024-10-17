package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AccessLogger *zap.Logger
	DBLogger     *zap.Logger
)

func InitLoggers() error {
	// Конфигурация для Access логов
	accessConfig := zap.NewProductionConfig()
	accessConfig.OutputPaths = []string{
		"access.log",
	}
	accessConfig.EncoderConfig.TimeKey = "timestamp"
	accessConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	AccessLogger, err = accessConfig.Build()
	if err != nil {
		return err
	}

	// Конфигурация для DB логов
	dbConfig := zap.NewProductionConfig()
	dbConfig.OutputPaths = []string{
		"db.log",
	}
	dbConfig.EncoderConfig.TimeKey = "timestamp"
	dbConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	DBLogger, err = dbConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

func SyncLoggers() {
	AccessLogger.Sync()
	DBLogger.Sync()
}
