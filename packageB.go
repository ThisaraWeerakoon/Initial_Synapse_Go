package main

import (
        "context"
        "log/slog"
        "sync"
        "sync/atomic"
		"os"

        "github.com/knadh/koanf/v2"
)

type CustomLoggerB struct {
        logger  *slog.Logger
        level   atomic.Int64 // Use atomic for safe concurrent access
        mu      sync.RWMutex // Mutex for safe access to handler options
        handler slog.Handler // Store the handler
}

var packageB *CustomLoggerB

func init() {
        packageB = newCustomLoggerB("packageB")
        packageB.setLevel(k) // Set initial level from config
}

func newCustomLoggerB(name string) *CustomLoggerB {
        // Default handler options
        opts := slog.HandlerOptions{
                Level: slog.LevelInfo, // Default level
        }
        handler := slog.NewTextHandler(os.Stdout, &opts)
        logger := slog.New(handler)

        return &CustomLoggerB{
                logger:  logger,
                handler: handler,
        }
}

func (l *CustomLoggerB) setLevel(k *koanf.Koanf) {
        levelStr := k.String("packageB") // Get the log level for packageB
        level := slog.LevelInfo         // Default level

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
		packageB.logger.Info("Package B log INFO")
		packageB.logger.Debug("Package B log DEBUG")
        // slog.Info("Package B log level updated", "level", level)
}


func main() {
        // Keep the main goroutine alive to allow watching.
        ctx := context.Background()
        <-ctx.Done() // Block indefinitely. You might want a cleaner way to exit in a real app.
}