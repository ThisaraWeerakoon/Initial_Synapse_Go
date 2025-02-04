package utils

import(
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

//This need to be changed later. This type of conversion is completely a result of our mock implementation of the core
func ProcessedMessageConverter(extractedFileData models.ExtractedFileData,isSuccess bool) models.ProcessedMessage{
	return models.ProcessedMessage{
		FilePath:extractedFileData.FilePath,
		IsSuccess:isSuccess,
	}
		
}