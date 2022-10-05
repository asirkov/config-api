package log

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/middleware"
)

const Version = "1"
const Component = "config-api"

var config Config
var instance Logger

func Init(c Config, l Logger) {
	config = c
	instance = l
}

// Logs
func infoCtx(ctx context.Context, message string) {
	instance.Log(
		Entry{
			LogLevel:   InfoLevel,
			Message:    message,
			ThreadName: middleware.GetReqID(ctx),
			Version:    Version,
			Component:  Component,
		},
	)
}

func InfoCtx(ctx context.Context, message string) {
	infoCtx(ctx, message)
}

func Info(message string) {
	infoCtx(nil, message)
}

func Infof(format string, args ...interface{}) {
	infoCtx(nil, fmt.Sprintf(format, args...))
}

func warnCtx(ctx context.Context, message string, err error) {
	if err != nil {
		message = fmt.Sprintf("%s %s", message, err)
	}

	instance.Log(
		Entry{
			LogLevel:   WarnLevel,
			Message:    message,
			ThreadName: middleware.GetReqID(ctx),
			Version:    Version,
			Component:  Component,
		},
	)
}

func WarnCtx(ctx context.Context, message string, err error) {
	warnCtx(ctx, message, err)
}

func Warn(message string, err error) {
	warnCtx(nil, message, err)
}

func errorCtx(ctx context.Context, message string, err error) {
	if err != nil {
		message = fmt.Sprintf("%s %s", message, err)
	}

	instance.Log(
		Entry{
			LogLevel:   ErrorLevel,
			Message:    message,
			ThreadName: middleware.GetReqID(ctx),
			Version:    Version,
			Component:  Component,
		},
	)
}

func ErrorCtx(ctx context.Context, message string, err error) {
	errorCtx(ctx, message, err)
}

func Error(message string, err error) {
	errorCtx(nil, message, err)
}
