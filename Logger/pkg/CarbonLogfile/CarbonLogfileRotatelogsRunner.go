package carbonlogfile

import (
	"log"
)

func CarbonLogfileRotatelogsRuner() {
	logger, err := initCarbonLogfileLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Example usage:
	logger.Debug("This is a debug message.")
	logger.Info("This is an info message.")
}