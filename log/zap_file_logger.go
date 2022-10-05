package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapFileLogger struct {
	Logger

	localLog *zap.Logger
}

func NewZapFileLogger(config Config) (Logger, error) {
	return NewZapFileLoggerWithCallerSkip(config, 4)
}

func NewZapFileLoggerWithCallerSkip(config Config, callerSkip int) (Logger, error) {
	var zapLogger ZapFileLogger

	cfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "@timestamp",
		MessageKey:     "message",
		CallerKey:      "logger_name",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	if config.LogFilesDestination == "" {
		return nil, fmt.Errorf("missed logFilesDestination configuration property")
	}
	if config.MaxSize == 0 {
		return nil, fmt.Errorf("missed maxSize configuration property")
	}
	if config.MaxAge == 0 {
		return nil, fmt.Errorf("missed maxAge configuration property")
	}
	if config.MaxBackups == 0 {
		return nil, fmt.Errorf("missed maxBackups configuration property")
	}

	filesOut := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s%s-current.log", config.LogFilesDestination, Component),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	})
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	logLevel := getLogLevel(config.LogLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, filesOut, logLevel),
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	zapLogger.localLog = zap.New(core, options...)
	return &zapLogger, nil
}

func (zl *ZapFileLogger) Info(message string, fields []zap.Field) {
	zl.localLog.Info(message, fields...)
}

func (zl *ZapFileLogger) Warn(message string, fields []zap.Field) {
	zl.localLog.Warn(message, fields...)
}

func (zl *ZapFileLogger) Error(message string, fields []zap.Field) {
	zl.localLog.Error(message, fields...)
}

func (zl *ZapFileLogger) Log(entry Entry) {
	fields := []zap.Field{}

	if entry.ThreadName != "" {
		fields = append(fields, zap.String("thread_name", entry.ThreadName))
	}

	if entry.Component != "" {
		fields = append(fields, zap.String("component", entry.Component))
	}

	if entry.Version != "" {
		fields = append(fields, zap.String("version", entry.Version))
	}

	for k, v := range entry.ExtraFields {
		if v != nil {
			fields = append(fields, zap.Any(k, v))
		}
	}

	switch entry.LogLevel {
	case InfoLevel:
		zl.Info(entry.Message, fields)
		break
	case WarnLevel:
		zl.Warn(entry.Message, fields)
		break
	case ErrorLevel:
		zl.Error(entry.Message, fields)
		break
	}
}

func getLogLevel(level Level) zapcore.Level {
	switch level {
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
