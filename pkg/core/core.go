// Note this is a mock implementation of the core. For now, The core will be responsible for starting the adapters and stopping them.
package core

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/fileinboundadapter"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/utils"
)

type FileInboundAdapterInterface interface {
	//start polling
	Start()
	//reveive results
	//stop process
}

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ReceiveRequests(*models.ExtractedFileData) {
	//receive requests

}

//This is a mock implementation of the parsing. This should be implemented in the future
func (c *Core) MockParsing(input *models.ExtractedFileData) {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 99
	chance := rand.Intn(100) // Generates a number in [0, 99]

	// 20% chance of failure
	if chance < 20 {
		utils.ProcessedMessageConverter(*input, false)
	} else {
		utils.ProcessedMessageConverter(*input, true)
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

