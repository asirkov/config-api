package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapConsoleLogger struct {
	Logger

	localLog *zap.Logger
}

const levelFormat = "[%5.5s]"
const levelColoredFormat = "\x1b[%dm[%5.5s]\x1b[0m"
const timeFormat = "[%s]"
const timeLayout = "15:04:05"
const callerFormat = "[%30.30s]"

var (
	_levelToColor = map[zapcore.Level]uint8{
		zap.InfoLevel:  34,
		zap.WarnLevel:  33,
		zap.ErrorLevel: 31,
	}
)

func NewZapConsoleLogger(config Config) (Logger, error) {
	return NewZapConsoleLoggerWithCallerSkip(config, 4)
}

func NewZapConsoleLoggerWithCallerSkip(config Config, callerSkip int) (Logger, error) {
	var zapLogger ZapConsoleLogger

	encodeLevel := func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		c, ok := _levelToColor[l]
		if ok {
			enc.AppendString(fmt.Sprintf(levelColoredFormat, c, l.CapitalString()))
		} else {
			enc.AppendString(fmt.Sprintf(levelFormat, l.CapitalString()))
		}
	}

	encodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf(timeFormat, t.Format(timeLayout)))
	}

	encodeCaller := func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf(callerFormat, c.TrimmedPath()))
	}

	consoleCfg := zapcore.EncoderConfig{
		LevelKey:      "L",
		TimeKey:       "T",
		MessageKey:    "M",
		CallerKey:     "C",
		StacktraceKey: "",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   encodeLevel,
		EncodeTime:    encodeTime,
		EncodeCaller:  encodeCaller,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	logLevel := getLogLevel(config.LogLevel)

	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel))

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
	}

	zapLogger.localLog = zap.New(core, options...)
	return &zapLogger, nil
}

func (zl *ZapConsoleLogger) Info(message string) {
	zl.localLog.Info(message)
}

func (zl *ZapConsoleLogger) Warn(message string) {
	zl.localLog.Warn(message)
}

func (zl *ZapConsoleLogger) Error(message string) {
	zl.localLog.Error(message)
}

func (zl *ZapConsoleLogger) Log(entry Entry) {
	switch entry.LogLevel {
	case InfoLevel:
		zl.Info(entry.Message)
		break
	case WarnLevel:
		zl.Warn(entry.Message)
		break
	case ErrorLevel:
		zl.Error(entry.Message)
		break
	}
}
