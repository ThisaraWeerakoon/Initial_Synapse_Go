package mypackage

import (
    "context"
    "log/slog"
    "os"
)

var logger *slog.Logger

func init() {
    // Create a logger with a specific prefix for this package
    logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,  // Optional: add source info
        Level:     slog.LevelDebug, // Package-specific level
    })).With("package", "mypackage") // Add package context
}

func MyFunction() {
    logger.Debug("This is a debug message from mypackage")
    logger.Info("This is an info message from mypackage")

    // Example of using context for more granular control within the package
    ctx := context.Background()
    ctx = context.WithValue(ctx, "operation", "MyFunction") // Add operation context
    logger.InfoContext(ctx, "Starting operation")
    // ...
}