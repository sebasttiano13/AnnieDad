package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type CTXValue string

const (
	ContextEventID   CTXValue = "event_id"
	ContextRequestID CTXValue = "request_id"
)

var contextValues = []CTXValue{ContextEventID, ContextRequestID}

func Debug(str string) {
	slog.Default().Debug(str)
}

func Info(str string) {
	slog.Default().Info(str)
}

func Warning(str string) {
	slog.Default().Warn(str)
}

func Error(str string) {
	slog.Default().Error(str)
}

func Debugf(fstring string, formaters ...any) {
	slog.Default().Debug(fmt.Sprintf(fstring, formaters...))
}

func Infof(fstring string, formaters ...any) {
	slog.Default().Info(fmt.Sprintf(fstring, formaters...))
}

func Warningf(fstring string, formaters ...any) {
	slog.Default().Warn(fmt.Sprintf(fstring, formaters...))
}

func Errorf(fstring string, formaters ...any) {
	slog.Default().Error(fmt.Sprintf(fstring, formaters...))
}

func DebugfContext(ctx context.Context, fstring string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Debug(fmt.Sprintf(fstring, formaters...))
}

func InfofContext(ctx context.Context, fstring string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Info(fmt.Sprintf(fstring, formaters...))
}

func WarnfContext(ctx context.Context, fstring string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Warn(fmt.Sprintf(fstring, formaters...))
}

func ErrorfContext(ctx context.Context, fstring string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Error(fmt.Sprintf(fstring, formaters...))
}

func addAttrsToLogger(ctx context.Context, log *slog.Logger) *slog.Logger {
	logger := log
	for _, val := range contextValues {
		if attr, ok := ctx.Value(val).(string); ok {
			logger = logger.With(slog.String(string(val), attr))
		}
	}

	return logger
}
