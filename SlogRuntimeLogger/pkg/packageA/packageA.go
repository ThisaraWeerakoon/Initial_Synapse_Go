package packageA

import (
	"SlogRuntimeLogger/pkg/logging"
	"SlogRuntimeLogger/pkg/packageC"
	"fmt"
)

type PackageA struct {
	User string
	Caller  string
	LoggerInstance logging.Logger
}

// New creates a new packageA.
func New(user string, caller string, loggerInstance logging.Logger) *PackageA {
	loggerInstance.PackageName = "packageA"
	
	return &PackageA{
		User: user,
		Caller:  caller,
		LoggerInstance: loggerInstance,
	}
}

func (p *PackageA) Run() {
	// fmt.Println("Logging level at package A is:", p.logger.Level())
	fmt.Println("Executing Run in package A...........................................................")
	p.LoggerInstance.Debugw("This is a debug message", "Package", p.User, "Caller", p.Caller)
	p.LoggerInstance.Infow("This is an info message", "Package", p.User, "Caller", p.Caller)
	p.LoggerInstance.Warnw("This is a warning message", "Package", p.User, "Caller", p.Caller)
	p.LoggerInstance.Errorw("This is an error message", "Package", p.User, "Caller", p.Caller)
	packageC := packageC.New("C", "A", p.LoggerInstance)
	packageC.Run()

}
