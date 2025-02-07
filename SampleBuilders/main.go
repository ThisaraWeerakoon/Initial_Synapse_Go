package main

import (
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/SampleBuilders/pkg/formurlencodedbuilder"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/SampleBuilders/pkg/jsonbuilder"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/SampleBuilders/pkg/multipartbuilder"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/SampleBuilders/pkg/xmlbuilder"
)

func main() {
	jsonbuilder.JSONBuilderRunner()
	xmlbuilder.XMLBuilderRunner()
	multipartbuilder.MultipartBuilderRunner()
	formurlencodedbuilder.FormUrlEncodedBuilderRunner()
	jsonbuilder.JSONBuilderRunner()

}
