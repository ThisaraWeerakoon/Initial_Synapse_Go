package core

func (c *Core) FileInboundAdapterRunner(fileInboundAdapterInterface FileInboundAdapterInterface) {
	go fileInboundAdapterInterface.Start()
}