package main

import (
        "context"
        "fmt"
        "log/slog"
        "os"
        "sync"
        "time"

        "gopkg.in/yaml.v3" // Or any other YAML library you prefer
)

// Config represents the structure of your log configuration file.
type Config struct {
        Packages map[string]slog.Level `yaml:"packages"`
}

// LogManager manages loggers for different packages and handles config changes.
type LogManager struct {
        configPath string
        mu         sync.RWMutex
        loggers    map[string]*slog.Logger
        config     *Config
        cancel     context.CancelFunc // For stopping the file watcher
}

func NewLogManager(configPath string) (*LogManager, error) {
        lm := &LogManager{
                configPath: configPath,
                loggers:    make(map[string]*slog.Logger),
        }

        if err := lm.loadConfig(); err != nil {
                return nil, fmt.Errorf("initial config load: %w", err)
        }

        if err := lm.setupLoggers(); err != nil {
                return nil, fmt.Errorf("initial logger setup: %w", err)
        }

        if err := lm.watchConfig(); err != nil {
                return nil, fmt.Errorf("starting config watcher: %w", err)
        }

        return lm, nil
}

func (lm *LogManager) loadConfig() error {
        lm.mu.Lock()
        defer lm.mu.Unlock()

        configFile, err := os.ReadFile(lm.configPath)
        if err != nil {
                return fmt.Errorf("reading config file: %w", err)
        }

        var config Config
        if err := yaml.Unmarshal(configFile, &config); err != nil {
                return fmt.Errorf("unmarshaling config: %w", err)
        }

        lm.config = &config
        return nil
}

func (lm *LogManager) setupLoggers() error {
        lm.mu.Lock()
        defer lm.mu.Unlock()

        for pkg, level := range lm.config.Packages {
                lm.loggers[pkg] = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
        }
        return nil
}

func (lm *LogManager) watchConfig() error {
        ctx, cancel := context.WithCancel(context.Background())
        lm.cancel = cancel

    // Use a file system watcher library or implement polling.  This example uses polling for simplicity.
        go func() {
                lastModTime := time.Time{}
                for {
                        select {
                        case <-ctx.Done():
                                return
                        case <-time.After(time.Second): // Check every second (adjust as needed)
                                fileInfo, err := os.Stat(lm.configPath)
                                if err != nil {
                                        slog.Error("watching config file", "error", err)
                                        continue
                                }

                                modTime := fileInfo.ModTime()
                                if modTime.After(lastModTime) {
                                        slog.Info("config file changed, reloading")
                                        if err := lm.loadConfig(); err != nil {
                                                slog.Error("reloading config", "error", err)
                                                continue
                                        }
                                        if err := lm.setupLoggers(); err != nil {
                                                slog.Error("resetting loggers", "error", err)
                                                continue
                                        }
                    lastModTime = modTime // Important: Update lastModTime *after* successful reload.
                                }
                        }
                }
        }()
        return nil
}

func (lm *LogManager) Logger(pkg string) *slog.Logger {
        lm.mu.RLock()
        defer lm.mu.RUnlock()

        if logger, ok := lm.loggers[pkg]; ok {
                return logger
        }
        // Provide a default logger or handle the case where the package is not configured.
        return slog.Default() // Or create a no-op logger.
}

func (lm *LogManager) Close() {
        if lm.cancel != nil {
                lm.cancel()
        }
}


func main() {
        configPath := "config.yaml" // Path to your YAML configuration file.

//     // create dummy config file
//     createDummyConfigFile(configPath)

        lm, err := NewLogManager(configPath)
        if err != nil {
                panic(err)
        }
        defer lm.Close()

        // Example usage:
        loggerA := lm.Logger("packageA")
        loggerB := lm.Logger("packageB")
        loggerC := lm.Logger("packageC")

        for {
                loggerA.Debug("This is a debug message from package A")
                loggerB.Info("This is an info message from package B")
                loggerC.Warn("This is a warning message from package C")
                time.Sleep(5*time.Second)
        }



}

func createDummyConfigFile(configPath string) {
    config := Config{
        Packages: map[string]slog.Level{
            "packageA": slog.LevelDebug,
            "packageB": slog.LevelInfo,
            "packageC": slog.LevelWarn,
        },
    }
    configFile, _ := yaml.Marshal(config)
    os.WriteFile(configPath, configFile, 0644)
}

func changeLogLevel(configPath string, packageName string, newLevel slog.Level) {
    configFile, _ := os.ReadFile(configPath)
    var config Config
    yaml.Unmarshal(configFile, &config)
    config.Packages[packageName] = newLevel
    newConfigFile, _ := yaml.Marshal(config)
    os.WriteFile(configPath, newConfigFile, 0644)
}
    