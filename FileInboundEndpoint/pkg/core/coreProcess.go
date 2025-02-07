package core


import (
	"math/rand"
	"time"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/FileInboundEndpoint/pkg/models"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/FileInboundEndpoint/pkg/utils"
)

//After receiving the data from file adapter initialize the CoreProcess which completes up to sending the results to the core. 
// Attention : Here need to make this function access through core interface should become handle concurrent calling from readingFile threads.
func (c *Core) ReceiveRequests(extractedFileDataFromFileAdapter *models.ExtractedFileDataFromFileAdapter){
	c.CoreProcess(extractedFileDataFromFileAdapter)

}

//This is a mock implementation of the parsing. This should be implemented in the future
func (c *Core) MockParsing(input *models.ExtractedFileDataFromFileAdapter) models.ProcessedMessageFromCore{
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 99
	chance := rand.Intn(100) // Generates a number in [0, 99]

	// 20% chance of failure
	if chance < 20 {
		return utils.ProcessedMessageConverter(*input, false)
	} else {
		return utils.ProcessedMessageConverter(*input, true)
	}

}

//this is for sequential steps after receving the message to core from sending the processed results to the adapter
func (c *Core) CoreProcess (extractedFileDataFromFileAdapter *models.ExtractedFileDataFromFileAdapter){
	//start mock parsing
	processedMessageFromCore := c.MockParsing(extractedFileDataFromFileAdapter)

	//send the results to the fileinboundadapter
	c.FileInboundAdapter.ReceiveResults(
		processedMessageFromCore,
	)
}