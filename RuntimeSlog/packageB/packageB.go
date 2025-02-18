package packageB

import (
	// "log/slog"
	// "os"

	"github.com/knadh/koanf/v2"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/models"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/logger"
)

var packageBLogger *models.CustomLogger

func init() {
	packageBLogger = logger.InitializeLogger("packageB")
	logger.RegisterLogger("packageB", packageBSetLevel{}) // Register
}

type packageBSetLevel struct{}

func (p packageBSetLevel) SetLevel(k *koanf.Koanf, name string) {
    packageBLogger.SetLevel(k, name)
}

func PackageBFunction() {
	packageBLogger.Logger.Info("Package B log INFO")
	packageBLogger.Logger.Debug("Package B log DEBUG")
}
