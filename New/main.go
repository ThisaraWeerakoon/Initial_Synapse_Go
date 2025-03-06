package main

import (
	"New/pkg/config"
	"New/pkg/packageA"
	"New/pkg/packageB"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	configFilePath := "config.yaml"
	cfg, err := config.InitializeConfig(configFilePath)
	if err != nil {
		log.Fatalf("Initialization error: %s", err.Error())
	}

	// Create package instances - they will auto-register when getting loggers
	packageA := packageA.New("A", "main")
	packageB := packageB.New("B", "main")

	// Start watching for config changes
	cfg.Watch(context.Background(), configFilePath)

	for {
		fmt.Println("Looping............................................................")
		
		packageA.InitializePackageC()
		packageA.Run()
		packageB.Run()
		time.Sleep(5 * time.Second)
	}
}
