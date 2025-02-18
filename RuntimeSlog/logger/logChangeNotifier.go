package logger

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/packageA"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/packageB"
	"github.com/knadh/koanf/v2"
)




func notifier(dataChan <-chan *koanf.Koanf) {
	for newK := range koanfChan {
		if newK.String("packageA") != currentK.String("packageA") {
			packageA.SetLevel(newK)
		}

		if newK.String("packageB") != currentK.String("packageB") {
			packageB.SetLevel(newK)
		}


	}

	// This is a dummy function to make the package a valid Go package.
}