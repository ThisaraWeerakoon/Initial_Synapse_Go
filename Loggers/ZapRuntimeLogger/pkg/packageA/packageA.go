package packageA

import (
	"fmt"
	"ZapRuntimeLogger/pkg/packageB"

	"go.uber.org/zap"
)


type PackageA struct {
	User string
	Age  int
	logger *zap.SugaredLogger
	
}

func New(user string, age int, loggerMap map[string]*zap.SugaredLogger) *PackageA {
	// logger = logger.Named("packageA")
	logger := loggerMap["packageA"]
	return &PackageA{
		User: user,
		Age:  age,
		logger: logger,
	}
}

func (p *PackageA) Run() {
	fmt.Println("Logging level at package A is:", p.logger.Level())
	p.logger.Debug("This is a debug message")
	p.logger.Info("This is an info message")
	packageB := packageB.New("Jane", 25, p.logger)
	packageB.Run()

}
