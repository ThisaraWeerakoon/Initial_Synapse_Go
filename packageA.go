package main

import (
	"log/slog"
	"os"
	"sync"
	"sync/atomic"

	"github.com/knadh/koanf/v2"
)

type CustomLogger struct {
	logger  *slog.Logger
	level   atomic.Int64 // Use atomic for safe concurrent access
	mu      sync.RWMutex // Mutex for safe access to handler options
	handler slog.Handler // Store the handler
}

var packageA *CustomLogger

func init() {
	packageA = newCustomLogger("packageA")
	packageA.setLevel(k) // Set initial level from config
}

func newCustomLogger(name string) *CustomLogger {
	// Default handler options
	opts := slog.HandlerOptions{
		Level: slog.LevelInfo, // Default level
	}
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(handler)

	return &CustomLogger{
		logger:  logger,
		handler: handler,
	}
}

func (l *CustomLogger) setLevel(k *koanf.Koanf) {
	levelStr := k.String("packageA") // Get the log level for packageA
	level := slog.LevelInfo          // Default level

	switch levelStr {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	l.mu.Lock()
	// Default handler options
	opts := slog.HandlerOptions{
		Level: level, // New level
	}
	l.handler = slog.NewTextHandler(os.Stdout, &opts)
	l.mu.Unlock()

	l.logger = slog.New(l.handler) // Create a new logger with the updated handler
	l.level.Store(int64(level))
	packageA.logger.Info("Package A log INFO")
	packageA.logger.Debug("Package A log DEBUG")
	//slog.Info("Package A log level updated", "level", level)
}
