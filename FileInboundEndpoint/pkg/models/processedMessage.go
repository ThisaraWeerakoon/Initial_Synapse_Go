package models

// Message sending from core to fileInboundAdapter representing the status of the previous message contexts that has being sent from adapter. 
// This need to be more organized including the error logs while parsing the message after completing the core in future
type ProcessedMessageFromCore struct {
	FilePath string
	IsSuccess bool
}