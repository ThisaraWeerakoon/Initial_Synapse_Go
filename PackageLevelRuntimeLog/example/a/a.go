package a

import "os"

var _, skipSync = os.LookupEnv("SKIP_SYNC")

func A() {
	logger := getLogger()
	if !skipSync {
		defer logger.Sync()
	}

	logger.Debug("debug message from package a")
	logger.Info("info message from package a")
	logger.Warn("warn message from package a")
	logger.Error("error message from package a")
}