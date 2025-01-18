package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CustomColorEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorStart, colorEnd string

	switch level {
	case zap.DebugLevel:
		colorStart = "\033[34m"
	case zap.InfoLevel:
		colorStart = "\033[32m"
	case zap.WarnLevel:
		colorStart = "\033[33m"
	case zap.ErrorLevel:
		colorStart = "\033[31m"
	case zap.DPanicLevel:
		colorStart = "\033[35m"
	default:
		colorStart = "\033[0m"
	}

	colorEnd = "\033[0m"
	enc.AppendString(fmt.Sprintf("%s%s%s", colorStart, level.CapitalString(), colorEnd))
}

var Logger *zap.Logger

func InitLogger() {

	logFileName := fmt.Sprintf("logs_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}

	writer := zapcore.AddSync(colorable.NewColorableStdout())
	fileWriter := zapcore.AddSync(logFile)

	multiWriter := zapcore.NewMultiWriteSyncer(writer, fileWriter)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    CustomColorEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		multiWriter,
		zapcore.DebugLevel,
	)

	Logger = zap.New(core, zap.WithCaller(true))
}
