package b

import "os"

var _, skipSync = os.LookupEnv("SKIP_SYNC")

func B() {
	logger := getLogger()
	if !skipSync {
		defer logger.Sync()
	}

	logger.Debug("debug message from package b")
	logger.Info("info message from package b")
	logger.Warn("warn message from package b")
	logger.Error("error message from package b")
}

