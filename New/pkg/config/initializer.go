package config

import (
	"New/pkg/loggerfactory"
	"context"
	"fmt"
)

func InitializeConfig(configFilePath string) (*Config, error) {
	cfg, err := ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var levelMap map[string]string
	var slogHandlerConfig loggerfactory.SlogHandlerConfig

	if cfg.IsSet("logger") {
		cfg.MustUnmarshal("logger.handler", &slogHandlerConfig)
		cfg.MustUnmarshal("logger.level.packages", &levelMap)
	}

	cm := loggerfactory.GetConfigManager()
	cm.SetLogLevelMap(&levelMap)
	cm.SetSlogHandlerConfig(slogHandlerConfig)

	cfg.Watch(context.Background(), configFilePath)

	return cfg, nil
}