package core

import(
	"context"
	"sync"
)

func (c *Core) FileInboundAdapterRunner(fileInboundAdapterInterface FileInboundAdapterInterface) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go fileInboundAdapterInterface.Start(ctx, &wg)
	wg.Wait()
}