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
	fileInboundAdapter FileInboundAdapterInterface
}

func NewCore(fileInboundAdapter FileInboundAdapterInterface) *Core {
	return &Core{
		fileInboundAdapter: fileInboundAdapter,
	}
}




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

	fileInboundAdapter := fileinboundadapter.NewFileInboundAdapter(config, c)
	c.FileInboundAdapterRunner(fileInboundAdapter) //start the fileinboundadapter. There may be another adapters run concurrenctly after the full implementation of the core

}

func (c *Core) Stop() {
	//stop the core

}

