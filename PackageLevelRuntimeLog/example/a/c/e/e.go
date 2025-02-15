package e

import "os"

var _, skipSync = os.LookupEnv("SKIP_SYNC")

func E() {
	logger := getLogger()
	if !skipSync {
		defer logger.Sync()
	}

	logger.Debug("debug message from package a/c/e")
	logger.Info("info message from package a/c/e")
	logger.Warn("warn message from package a/c/e")
	logger.Error("error message from package a/c/e")
}