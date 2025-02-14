package main

import (
	carbonconsolelogging "github.com/ThisaraWeerakoon/Initial_Synapse_Go/Logger/pkg/CarbonConsole"
	carbonlogfilezaplamberjack "github.com/ThisaraWeerakoon/Initial_Synapse_Go/Logger/pkg/CarbonLogfile"
)
func main() {
	// This is the main function for the Logger application.

	//CarbonConsoleLogging
	carbonconsolelogging.CarbonConsoleLoggingRunner()

	//CarbonLogfileZapLamberjack
	carbonlogfilezaplamberjack.CarbonLogfileZapLamberjackRuner()

}