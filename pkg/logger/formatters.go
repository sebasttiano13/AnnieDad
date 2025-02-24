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

func Debugf(fsting string, formaters ...any) {
	slog.Default().Debug(fmt.Sprintf(fsting, formaters...))
}

func Infof(fsting string, formaters ...any) {
	slog.Default().Info(fmt.Sprintf(fsting, formaters...))
}

func Warningf(fsting string, formaters ...any) {
	slog.Default().Warn(fmt.Sprintf(fsting, formaters...))
}

func Errorf(fsting string, formaters ...any) {
	slog.Default().Error(fmt.Sprintf(fsting, formaters...))
}

func DebugfContext(ctx context.Context, fsting string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Debug(fmt.Sprintf(fsting, formaters...))
}

func InfofContext(ctx context.Context, fsting string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Info(fmt.Sprintf(fsting, formaters...))
}

func WarnfContext(ctx context.Context, fsting string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Warn(fmt.Sprintf(fsting, formaters...))
}

func ErrorfContext(ctx context.Context, fsting string, formaters ...any) {
	logger := addAttrsToLogger(ctx, slog.Default())
	logger.Error(fmt.Sprintf(fsting, formaters...))
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
