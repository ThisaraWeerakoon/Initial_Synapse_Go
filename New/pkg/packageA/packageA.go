package packageA

import (
	"New/pkg/loggerfactory"
	"New/pkg/packageC"
	"fmt"
	"log/slog"
)

var(
	InstanceOfC *packageC.PackageC
)

type PackageA struct {
	User   string
	Caller string
	logger *slog.Logger
}

// New creates a new packageA.
func New(user string, caller string) *PackageA {
	p := &PackageA{
		User:   user,
		Caller: caller,
	}
	// Pass the instance to GetLogger for automatic registration
	p.logger = loggerfactory.GetLogger("packageA", p)
	return p
}

func (p *PackageA) UpdateLogger() {
	p.logger = loggerfactory.GetLogger("packageA", p)
}

func (p *PackageA) InitializePackageC (){
	InstanceOfC = packageC.New("C","A")
	
}

func (p *PackageA) Run() {
	fmt.Println("Executing Run in package A...........................................................")
	p.logger.Debug("This is a debug message", "Package", p.User, "Caller", p.Caller)
	p.logger.Info("This is an info message", "Package", p.User, "Caller", p.Caller)
	p.logger.Warn("This is a warning message", "Package", p.User, "Caller", p.Caller)
	p.logger.Error("This is an error message", "Package", p.User, "Caller", p.Caller)
	InstanceOfC.Run()
}
