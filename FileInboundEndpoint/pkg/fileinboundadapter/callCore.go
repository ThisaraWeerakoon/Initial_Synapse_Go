package fileinboundadapter

import(
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

//CallCore is a function that calls the core to receive the extracted file data from the file adapter.
//Initialize the process of the core by sending the extracted file data to the core.
func (f *FileInboundAdapter) CallCore(extractedFileDataFromFileAdapter models.ExtractedFileDataFromFileAdapter){
	f.core.ReceiveRequests(&extractedFileDataFromFileAdapter)


}