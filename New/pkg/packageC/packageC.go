package packageC

import (
	"New/pkg/loggerfactory"
	"fmt"
	"log/slog"
)

type PackageC struct {
	User string
	Caller  string
	logger *slog.Logger
}

// New creates a new packageC.
func New(user string, caller string) *PackageC {
	p := &PackageC{
		User: user,
		Caller:  caller,
	}
	p.logger = loggerfactory.GetLogger("packageC", p)
	return p
}

func (p *PackageC) UpdateLogger(){
	p.logger = loggerfactory.GetLogger("packageC",p)
}



func (p *PackageC) Run() {
	fmt.Println("Executing Run in package C...........................................................")
	p.logger.Debug("This is a debug message", "Package", p.User, "Caller", p.Caller)
	p.logger.Info("This is an info message", "Package", p.User, "Caller", p.Caller)
	p.logger.Warn("This is a warning message", "Package", p.User, "Caller", p.Caller)
	p.logger.Error("This is an error message", "Package", p.User, "Caller", p.Caller)

}