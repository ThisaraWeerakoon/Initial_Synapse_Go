package main

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/config"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/logger"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/packageA"
)

func main() {
	   //Initialize config
        cfg, err := config.New("config.toml")
        if err != nil {
            panic(err)
        }


		// //Initialize logger
    	foo := logger.NewMappedLogger(cfg)

		packageA.DoSomethingA()

        // a.DoSomethingA()
        // b.DoSomethingB()
        // c.DoSomethingC()

        // fmt.Println("Config file change detected. Log levels updated.")
        // select {}
}