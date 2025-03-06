package packageB

import (
	"New/pkg/loggerfactory"
	"fmt"
	"log/slog"
)

type PackageB struct {
	User string
	Caller  string
	logger *slog.Logger
}

// New creates a new packageA.
func New(user string, caller string) *PackageB {

	p := &PackageB{
		User: user,
		Caller:  caller,
	}
	p.logger = loggerfactory.GetLogger("packageB", p)
	return p
}

func (p *PackageB) UpdateLogger(){
	p.logger = loggerfactory.GetLogger("packageB", p)
}



func (p *PackageB) Run() {
	fmt.Println("Executing Run in package B...........................................................")
	p.logger.Debug("This is a debug message", "Package", p.User, "Caller", p.Caller)
	p.logger.Info("This is an info message", "Package", p.User, "Caller", p.Caller)
	p.logger.Warn("This is a warning message", "Package", p.User, "Caller", p.Caller)
	p.logger.Error("This is an error message", "Package", p.User, "Caller", p.Caller)

}