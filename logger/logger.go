package logger

import (
	"errors"

	"go.uber.org/zap"
)

func parseExtendedData(message string, extData []string) string {
	extendedData := ""
	for _, msg := range extData {
		extendedData += (" | " + msg)
	}

	return message + extendedData
}

func logger(loggerType, message, errorMessage string) {
	switch loggerType {
	case "info":
		Logger.Info(message)
	case "warn":
		Logger.Warn(message)
	case "debug":
		Logger.Debug(message, zap.Error(errors.New(errorMessage)))
	case "error":
		Logger.Error(message)
	case "panic":
		Logger.Panic(message, zap.Error(errors.New(errorMessage)))
	}
}

func Info(requestId string, message string, extData ...string) {
	loggerMessage := requestId + " | " + message
	logger("info", parseExtendedData(loggerMessage, extData), "")
}

func Error(requestId string, message string, errorMessage string, extData ...string) {
	loggerMessage := requestId + " | " + message + " | " + errorMessage
	logger("error", parseExtendedData(loggerMessage, extData), "")
}

func Warn(requestId string, message string, errorMessage string, extData ...string) {
	loggerMessage := requestId + " | " + message + " | " + errorMessage
	logger("warn", parseExtendedData(loggerMessage, extData), "")
}

func Debug(requestId string, message string, errorMessage string, extData ...string) {
	loggerMessage := requestId + " | " + message + " | " + errorMessage
	logger("debug", parseExtendedData(loggerMessage, extData), "")
}

func Panic(requestId string, message string, errorMessage string, extData ...string) {
	loggerMessage := requestId + " | " + message + " | " + errorMessage
	logger("panic", parseExtendedData(loggerMessage, extData), "")
}
