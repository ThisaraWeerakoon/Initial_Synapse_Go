package main

import (
	"ZapRuntimeLogger/pkg/config"
	"ZapRuntimeLogger/pkg/packageA"
	"ZapRuntimeLogger/pkg/packageB"
	"context"
	"fmt"
	"log"
	"time"

	"ZapRuntimeLogger/pkg/logging"

	"go.uber.org/zap"
)


func main() {
	configFilePath := "config.yaml"
	cfg, err := config.ReadFile(configFilePath)
	if err != nil {
			log.Fatalf("Cannot read config file: %s", err.Error())
	}
	var loggerConfig *zap.Config
	if cfg.IsSet("logger") {
		loggerConfig = &zap.Config{}
		cfg.MustUnmarshal("logger", loggerConfig)
	}
	logger, err := logging.NewNamedFromConfig(loggerConfig, "main")
	if err != nil {
		log.Fatalf("Cannot create logger: %s", err.Error())
	}

	packageA := packageA.New("John", 30, logger)
	packageB := packageB.New("Jane", 25, logger)

	cfg.Watch(context.Background(),configFilePath,loggerConfig)
	for {
		fmt.Println("Logging level at main is:", logger.Level())
		logger.Debug("This is a debug message")
		logger.Info("This is an info message")
		// logger.Warn("This is a warning message")
		// logger.Error("This is an error message")
		packageA.Run()
		packageB.Run()
		time.Sleep(5 * time.Second)
	}
}