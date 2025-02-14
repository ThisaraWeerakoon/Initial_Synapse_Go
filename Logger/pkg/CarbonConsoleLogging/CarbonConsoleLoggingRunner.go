package carbonconsolelogging

import (
	"github.com/sirupsen/logrus"
)

func CarbonConsoleLoggingRunner() {
	initConsoleLogger()

	// Example usage: adding a "logger" field to simulate category information.
	consoleLogger.WithField("logger", "org.example.MyComponent").
		Debug("This is a debug message.")

	// Example logging an error along with a message.
	consoleLogger.WithFields(logrus.Fields{
		"logger": "org.example.MyComponent",
		"error":  "file not found",
	}).Error("An error occurred")
}