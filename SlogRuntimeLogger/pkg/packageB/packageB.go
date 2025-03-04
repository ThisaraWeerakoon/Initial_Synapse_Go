package packageB

import (
	"SlogRuntimeLogger/pkg/logging"
	"fmt"
)

type PackageB struct {
	User string
	Caller  string
	LoggerInstance logging.Logger
}

// New creates a new packageA.
func New(user string, caller string, loggerInstance logging.Logger) *PackageB {
	loggerInstance.PackageName = "packageB"
	return &PackageB{
		User: user,
		Caller:  caller,
		LoggerInstance: loggerInstance,
	}
}

func (p *PackageB) Run() {
	// fmt.Println("Logging level at package A is:", p.logger.Level())
	fmt.Println("Executing Run in package B...........................................................")
	p.LoggerInstance.Debugw("This is a debug message", "User", p.User, "Caller", p.Caller)
	p.LoggerInstance.Infow("This is an info message", "User", p.User, "Caller", p.Caller)
	p.LoggerInstance.Warnw("This is a warning message", "User", p.User, "Caller", p.Caller)
	p.LoggerInstance.Errorw("This is an error message", "User", p.User, "Caller", p.Caller)

}
