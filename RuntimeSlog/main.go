package main

import (
	"fmt"
	"time"

	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/packageA"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/packageB"
)

func main() {
	for{
		fmt.Println("...............Starting the loop...............")
		packageA.PackageAFunction()
		packageB.PackageBFunction()
		time.Sleep(5*time.Second)
	}
}