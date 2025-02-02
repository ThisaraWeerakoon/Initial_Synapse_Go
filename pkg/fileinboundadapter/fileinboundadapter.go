package fileinboundadapter

type CoreInterface interface {
	//ReceiveRequests
}

type FileInboundAdapter struct{

}

func NewFileInboundAdapter() *FileInboundAdapter {
	return &FileInboundAdapter{}
}

func (f  *FileInboundAdapter) StartPolling() {
	//start polling
}

func (f  *FileInboundAdapter) ReceiveResults() {
	//reveive results
}