package packageA

import (
	"New/pkg/loggerfactory"
	"fmt"
	"log/slog"
	"New/pkg/packageC"
)

type PackageA struct {
	User string
	Caller  string
	logger *slog.Logger
}

// New creates a new packageA.
func New(user string, caller string) *PackageA {
	
	
	return &PackageA{
		User: user,
		Caller:  caller,
		logger: loggerfactory.GetLogger("packageA"),
	}

}

func (p *PackageA) UpdateLogger(){
	p.logger = loggerfactory.GetLogger("packageA")
}



func (p *PackageA) Run() {
	// fmt.Println("Logging level at package A is:", p.logger.Level())
	fmt.Println("Executing Run in package A...........................................................")
	p.logger.Debug("This is a debug message", "Package", p.User, "Caller", p.Caller)
	p.logger.Info("This is an info message", "Package", p.User, "Caller", p.Caller)
	p.logger.Warn("This is a warning message", "Package", p.User, "Caller", p.Caller)
	p.logger.Error("This is an error message", "Package", p.User, "Caller", p.Caller)
	packageC := packageC.New("C", "A")
	packageC.Run()

}