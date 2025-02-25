package logger

import (
	"fmt"
	"log/slog"
	"os"

	"runtime"
	"strings"
	"sync"

	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/config"
)

type RootLogger struct {

        config *config.Config

        mu sync.RWMutex

        //For now we only assume that only changing in runtime is log level.
        logLevelMap map[string]slog.Level // Package-specific log levels

        //There should be another attributes related to logging configurations like appender, file rotation condiftions etc.

}



var rootLoggerInstance *RootLogger
var once sync.Once


// InitializeLogger MUST be called once at the start of your application,
// passing in the config.  This sets up the singleton logger instance.
func InitializeLogger(cfg *config.Config) {
    once.Do(func() {
        rootLoggerInstance = New(cfg)
    })
}

// GetLogger returns the package level logger instance.  It panics if 
// InitializeLogger hasn't been called yet.
func GetLogger(packageName string) *slog.Logger {
    if rootLoggerInstance == nil {
        panic("logger not initialized")
    }

    opts := slog.HandlerOptions{
		Level: rootLoggerInstance.logLevelMap[packageName], // New level
	}

    //Here NewTextHandler is used. But in a full implementation, this should be configurable.
    handler := slog.NewTextHandler(os.Stdout, &opts)
    logger := slog.New(handler)
    return logger
}

func New(config *config.Config) *RootLogger {
    l := &RootLogger{
            config: config,
            logLevelMap: make(map[string]slog.Level),
    }

    config.AddObserver(l)
    l.Update(config.All()) // Initialize from initial config

    return l
}

func (l *RootLogger) Update(config map[string]interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()

    for key, value:= range config {
        if strings.HasPrefix(key, "log_levels.") {
            pkg:= strings.TrimPrefix(key, "log_levels.")
            levelString, ok:= value.(string)
            if!ok {
                fmt.Println("Invalid log level for package", pkg)
                continue
            }
            slogLevel:= parseLogLevel(levelString)
            l.logLevelMap[pkg] = slogLevel
        }
    }

    // You might want to reconfigure slog handlers here if needed
}

func parseLogLevel(level string) slog.Level {
    switch level {
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo // Default level
    }
}



// func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
//         l.logWithLevel(slog.LevelInfo, msg, args...)
// }

// func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
//         l.logWithLevel(slog.LevelDebug, msg, args...)
// }

// func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
//     l.logWithLevel(slog.LevelError, msg, args...)
// }

// func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
//     l.logWithLevel(slog.LevelWarn, msg, args...)
// }

// func (l *Logger) logWithLevel(level slog.Level, msg string, args ...interface{}) {
//     pkg := getPackageName() // Get package name using reflection 
//     l.mu.RLock()
//     defer l.mu.RUnlock()

//     if pkgLevel, ok := l.logLevels[pkg]; ok && level >= pkgLevel {
//         l.log.Log(context.Background(), level, msg, args...)
//     }
// }

// getPackageName uses reflection to get the calling package's name.
func getPackageName() string {
    pc := make([]uintptr, 1)
    runtime.Callers(2, pc) // Skip runtime.Callers and logWithLevel
    f := runtime.FuncForPC(pc[0])
    if f == nil {
            return "unknown"
    }
    fullName := f.Name()
    parts := strings.Split(fullName, ".")
    if len(parts) >= 2 {
        return strings.Join(parts[:len(parts)-1], ".") // Return package path
    }
    return "main" // Or a suitable default
}

