package b

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/pkglog"

	"go.uber.org/zap"
)

var logger *zap.Logger

func getLogger() *zap.Logger {
	if logger != nil {
		return logger
	}
	logger = pkglog.NewProductionLogger("github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/b")
	return logger
}

func ResetLogger() {
	logger = nil
}