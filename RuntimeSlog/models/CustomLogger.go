package models

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"os"


	"github.com/knadh/koanf/v2"
)

type CustomLogger struct {
	Logger  *slog.Logger
	Level   atomic.Int64 // Use atomic for safe concurrent access
	Mu      sync.RWMutex // Mutex for safe access to handler options
	Handler slog.Handler // Store the handler
}

// This function sets the koanf object into the logger. In full implementation this function will change the whole logging configuration. For now it sets only the log level for simplicity.
func (l *CustomLogger) SetLevel(k *koanf.Koanf, name string) {
	levelStr := k.String(name) // Get the log level for packageA
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

	l.Mu.Lock()
	// Default handler options
	opts := slog.HandlerOptions{
		Level: level, // New level
	}
	l.Handler = slog.NewTextHandler(os.Stdout, &opts)
	l.Mu.Unlock()

	l.Logger = slog.New(l.Handler) // Create a new logger with the updated handler
	l.Level.Store(int64(level))
}