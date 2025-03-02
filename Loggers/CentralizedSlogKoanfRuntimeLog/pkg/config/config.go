package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
        k *koanf.Koanf
        mu sync.RWMutex
        observers []Observer
}

type Observer interface {
        Update(config map[string]interface{})
}


func New(configPath string) (*Config, error) {
        k := koanf.New(".")

        // Load config file.
        if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
                return nil, fmt.Errorf("error loading config: %w", err)
        }

        c := &Config{k: k}
        if err := c.watchConfigChanges(configPath); err != nil {
                return nil, err
        }

        return c, nil
}

func (c *Config) Get(path string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.k.Get(path)
}

func (c *Config) All() map[string]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.k.All()
}

func (c *Config) AddObserver(o Observer) {
        c.mu.Lock()
        defer c.mu.Unlock()
        c.observers = append(c.observers, o)
}

func (c *Config) watchConfigChanges(configPath string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()

    err = watcher.Add(configPath)
    if err != nil {
        return err
    }

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    if err := c.reloadConfig(configPath); err != nil {
                        fmt.Println("Error reloading config:", err) // Or log using your logger
                    }
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                fmt.Println("error:", err) // Or log using your logger
            }
        }
    }()
    return nil
}

func (c *Config) reloadConfig(configPath string) error {
    c.mu.Lock() // Lock for writing
    defer c.mu.Unlock()

    newKoanf := koanf.New(".")
    if err := newKoanf.Load(file.Provider(configPath), toml.Parser()); err != nil {
        return fmt.Errorf("error reloading config: %w", err)
    }

    c.k = newKoanf // Atomically update the config

    // Notify observers in a separate goroutine
    go func() {
        configData := c.All()
        for _, o := range c.observers {
            o.Update(configData)
        }
    }()

    return nil
}
