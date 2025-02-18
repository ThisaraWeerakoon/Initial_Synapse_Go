package packageA

import (
	// "log/slog"
	// "os"

	"github.com/knadh/koanf/v2"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/models"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/logger"
)

var packageALogger *models.CustomLogger

func init() {
	packageALogger = logger.InitializeLogger("packageA")
}

func SetLevel(k *koanf.Koanf) {
	packageALogger.SetLevel(k, "packageA")
}

func PackageAFunction() {
	packageALogger.Logger.Info("Package A log INFO")
	packageALogger.Logger.Debug("Package A log DEBUG")
}
