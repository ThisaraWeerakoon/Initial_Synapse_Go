package c

import "os"

var _, skipSync = os.LookupEnv("SKIP_SYNC")

func C() {
	logger := getLogger()
	if !skipSync {
		defer logger.Sync()
	}

	logger.Debug("debug message from package a/c")
	logger.Info("info message from package a/c")
	logger.Warn("warn message from package a/c")
	logger.Error("error message from package a/c")
}