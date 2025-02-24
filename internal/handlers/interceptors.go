package handlers

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"log/slog"
)

// InterceptorLogger logging all incoming requests
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		attrs := make([]any, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok || i+1 >= len(fields) {
				continue
			}

			attrs = append(attrs, slog.Any(key, fields[i+1]))
		}

		logger := l.With(attrs...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			logger.Info(msg)
		}
	})
}
