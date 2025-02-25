package packageA

import (

	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/logger"
)
var log = logger.GetLogger("packageA")

func DoSomethingA() {
	
	log.Debug("This is a debug message from package A")  // No context needed!
	log.Info("This is an info message from package A")
	log.Error("This is an error message from package A")
}

