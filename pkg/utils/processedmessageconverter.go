package utils

import(
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

//This need to be changed later. This type of conversion is completely a result of our mock implementation of the core
func ProcessedMessageConverter(extractedFileData models.ExtractedFileDataFromFileAdapter,isSuccess bool) models.ProcessedMessageFromCore{
	return models.ProcessedMessageFromCore{
		FilePath:extractedFileData.FILE_PATH,
		IsSuccess:isSuccess,
	}
		
}