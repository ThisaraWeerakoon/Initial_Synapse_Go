package main

import (
	// "github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/config"
	// "github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/logger"
	"time"

	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/CentralizedSlogKoanfRuntimeLog/pkg/packageA"
)

func main() {
	//    //Initialize config
    //     cfg, err := config.New("config.toml")
    //     if err != nil {
    //         panic(err)
    //     }


		// //Initialize logger
        // logger.InitializeLogger(cfg)

    	// foo := logger.NewMappedLogger(cfg)
        for{
            packageA.DoSomethingA()
            time.Sleep(10 * time.Second)
        }


        // a.DoSomethingA()
        // b.DoSomethingB()
        // c.DoSomethingC()

        // fmt.Println("Config file change detected. Log levels updated.")
        // select {}
}