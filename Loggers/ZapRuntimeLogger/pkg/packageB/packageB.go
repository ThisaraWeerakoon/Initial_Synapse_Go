package packageB

import (
	"fmt"

	"go.uber.org/zap"
)


type PackageB struct {
	User string
	Age  int
	logger *zap.SugaredLogger
}

func New(user string, age int, logger *zap.SugaredLogger) *PackageB {
	logger = logger.Named("packageB")
	return &PackageB{
		User: user,
		Age:  age,
		logger: logger,
	}
}

func (p *PackageB) Run() {
	fmt.Println("Logging level at package B is:", p.logger.Level())
	p.logger.Debug("This is a debug message")
	p.logger.Info("This is an info message")

}
