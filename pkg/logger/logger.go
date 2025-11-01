package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger returns a singleton zap.Logger instance
func GetLogger() *zap.Logger {
	once.Do(func() {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

		// File encoder
		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

		// Log file
		logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fileWriter := zapcore.AddSync(logFile)
		consoleWriter := zapcore.AddSync(os.Stdout)

		// Set log level
		level := zapcore.InfoLevel

		// Combine cores
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleWriter, level),
			zapcore.NewCore(fileEncoder, fileWriter, level),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
	return logger
}

// Sugar returns a zap.SugaredLogger for convenience
func Sugar() *zap.SugaredLogger {
	return GetLogger().Sugar()
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}
