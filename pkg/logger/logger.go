package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	slogmulti "github.com/samber/slog-multi"
)

func marshalLevel(level string) (slog.Level, error) {
	if level == "" {
		return slog.LevelDebug, nil
	}
	var l slog.Level
	err := l.UnmarshalText([]byte(level))

	return l, err
}

func NewLogger(logFile string, level string) (*slog.Logger, error) {
	parsedLevel, err := marshalLevel(level)
	slogHandlers := make([]slog.Handler, 0, 2)
	if err != nil {
		return nil, err
	}
	var jsonHandler *slog.JSONHandler
	if logFile != "" {
		logFile, err := os.OpenFile(filepath.Clean(logFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize log file %w", err)
		}
		jsonHandler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: parsedLevel})
	} else {
		jsonHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parsedLevel})
	}

	slogHandlers = append(slogHandlers, jsonHandler)
	logger := slog.New(slogmulti.Fanout(slogHandlers...))

	slog.SetDefault(logger)

	return logger, nil
}

func NewTextLogger(level string) (*slog.Logger, error) {
	parsedLevel, err := marshalLevel(level)
	if err != nil {
		return nil, err
	}

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: parsedLevel})
	logger := slog.New(textHandler)
	slog.SetDefault(logger)

	return logger, nil
}

func GetDefault() *slog.Logger {
	return slog.Default()
}
