package core

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

type fileInboundAdapterInterface interface {
	//start polling
	//reveive results
	//stop process
}

type Core struct {
	Configurations //not generic have to change

}

func newCore() *Core {
	return &Core{}
}



func (c *Core) Start() {
	//start the core
}


