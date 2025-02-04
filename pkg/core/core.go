// Note this is a mock implementation of the core. For now, The core will be responsible for starting the adapters and stopping them.
package core

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/fileinboundadapter"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

type FileInboundAdapterInterface interface {
	//start polling
	Start()
	//reveive results
	ReceiveResults(processedMessageFromCore models.ProcessedMessageFromCore)
	//stop process
}

type Core struct {
	FileInboundAdapter FileInboundAdapterInterface
}

func NewCore() *Core {
	return &Core{}
}

// func (c *Core) InitializeFileInboundAdapter() {
// 	config := models.Configurations{
// 		Interval:           10,
// 		Sequential:         false,
// 		Coordination:       true,
// 		ActionAfterProcess: "MOVE",
// 		MoveAfterProcess:   "file:///home/thisarar/user/test/out",
// 		FileURI:            "file:///home/thisarar/user/test/in",
// 		MoveAfterFailure:   "file:///home/thisarar/user/test/failed",
// 		FileNamePattern:    "*.xml",
// 		ContentType:        "text/xml",
// 		ActionAfterFailure: "MOVE",
// 	}

// 	c.FileInboundAdapter = fileinboundadapter.NewFileInboundAdapter(config, c)
// }

func (c *Core) Run() {

	//start the fileinboundadapter. But this start should be more generic and user need to configure the configurations.Here I hardcoded the configurations
	config := models.Configurations{
		Interval:           10,
		Sequential:         false,
		Coordination:       true,
		ActionAfterProcess: "MOVE",
		MoveAfterProcess:   "file:///home/thisarar/user/test/out",
		FileURI:            "file:///home/thisarar/user/test/in",
		MoveAfterFailure:   "file:///home/thisarar/user/test/failed",
		FileNamePattern:    "*.xml",
		ContentType:        "text/xml",
		ActionAfterFailure: "MOVE",
	}

	//initialize the fileinboundadapter to core
	c.FileInboundAdapter = fileinboundadapter.NewFileInboundAdapter(config, c)
	c.FileInboundAdapterRunner(c.FileInboundAdapter) //start the fileinboundadapter. There may be another adapters run concurrenctly after the full implementation of the core

}

func (c *Core) Stop() {
	//stop the core

}

