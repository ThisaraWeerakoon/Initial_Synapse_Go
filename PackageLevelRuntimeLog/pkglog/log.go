package pkglog

import (
	"os"

	"go.uber.org/zap"
	"moul.io/zapfilter"
)

var filterRules zapfilter.FilterFunc

func Init(config string) {
	filterRules = zapfilter.MustParseRules(config)
}

func InitFromFile(configFile string) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	filterRules = zapfilter.MustParseRules(string(file))
}

func NewProductionLogger(namespace string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"output.log"}
	config.Sampling = nil
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	underlying, err := config.Build()
	if err != nil {
		panic(err)
	}

	return zap.New(zapfilter.NewFilteringCore(underlying.Core(), filterRules)).Named(namespace)
}

func NewDevelopmentLogger(namespace string) *zap.Logger {
	underlying, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return zap.New(zapfilter.NewFilteringCore(underlying.Core(), filterRules)).Named(namespace)
}