package d

import "os"

var _, skipSync = os.LookupEnv("SKIP_SYNC")

func D() {
	logger := getLogger()
	if !skipSync {
		defer logger.Sync()
	}

	logger.Debug("debug message from package b/d")
	logger.Info("info message from package b/d")
	logger.Warn("warn message from package b/d")
	logger.Error("error message from package b/d")
}