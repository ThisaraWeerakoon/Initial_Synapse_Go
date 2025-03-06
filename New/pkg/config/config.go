package config

import (
	"context"
	"fmt"
	"log"

	"New/pkg/loggerfactory"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	koanf *koanf.Koanf
}

func ReadFile(filename string) (*Config, error) {
	k := koanf.New(".")
	f := file.Provider(filename)
	if err := k.Load(f, yaml.Parser()); err != nil {
		return nil, err
	}

	cfg := &Config{
		koanf: k,
	}

	return cfg, nil
}

func (c *Config) IsSet(key string) bool {
	return c.koanf.Exists(key)
}

func (c *Config) Watch(ctx context.Context, filename string) {
	f := file.Provider(filename)

	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}
		// Throw away the old config and load a fresh copy.
		log.Println("config changed. Reloading ...")
		new_k := koanf.New(".")
		if err := new_k.Load(f, yaml.Parser()); err != nil {
			log.Printf("error loading new config: %v", err)
			return
		}
		// Update the config
		c.koanf = new_k

		// Update the logger configuration
		var levelMap map[string]string
		var slogHandlerConfig loggerfactory.SlogHandlerConfig

		c.MustUnmarshal("logger.level.packages", &levelMap)
		c.MustUnmarshal("logger.handler", &slogHandlerConfig)

		cm := loggerfactory.GetConfigManager()
		cm.SetLogLevelMap(&levelMap)
		cm.SetSlogHandlerConfig(slogHandlerConfig)
	})
}

func (c *Config) Unmarshal(key string, out interface{}) error {
	err := c.koanf.Unmarshal(key, out)
	if err != nil {
		return fmt.Errorf("cannot unmarshal config for key %q: %v", key, err)
	}
	return nil
}

func (c *Config) MustUnmarshal(key string, out interface{}) {
	err := c.Unmarshal(key, out)
	if err != nil {
		panic(err)
	}
}
