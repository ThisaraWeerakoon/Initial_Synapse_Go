package main

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/FileInboundEndpoint/pkg/core"

)



func main () {
	core := core.NewCore()
	core.Run()
}

