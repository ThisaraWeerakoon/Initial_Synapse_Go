package main

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/a"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/b"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/a/c"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/b/d"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/example/a/c/e"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelRuntimeLog/pkglog"
)

func init() {
	pkglog.InitFromFile("logging_props.txt")
}

func main() {
	a.A()
	b.B()
	c.C()
	d.D()
	e.E()
}