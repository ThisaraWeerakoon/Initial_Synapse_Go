package carbonconsole

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter implements logrus.Formatter to mimic the log4j pattern:
// Pattern: "[%d] %5p {%c{1}} - %m%ex%n"
type CustomFormatter struct{}

// Format builds the log entry string.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format timestamp similar to %d in log4j
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	
	// Format level (%5p) ensuring it's 5 characters wide, padded if necessary.
	level := strings.ToUpper(entry.Level.String())
	if len(level) < 5 {
		level = fmt.Sprintf("%-5s", level)
	}
	
	// Retrieve the logger category from entry.Data if available; otherwise use a default.
	// In log4j, %c{1} prints the last element of the logger name.
	category := "default"
	if cat, ok := entry.Data["logger"]; ok {
		// Here, we assume category is a string; in a real app you might split by '.' and take the last segment.
		catStr := fmt.Sprintf("%v", cat)
		parts := strings.Split(catStr, ".")
		category = parts[len(parts)-1]
	}
	
	// Log message (%m)
	message := entry.Message
	
	// Check if an error is included in the entry and append it similar to %ex in log4j.
	errStr := ""
	if errVal, ok := entry.Data["error"]; ok {
		errStr = fmt.Sprintf(" %v", errVal)
	}
	
	// Construct the log line matching the pattern "[timestamp] LEVEL {category} - message error\n"
	logLine := fmt.Sprintf("[%s] %s {%s} - %s%s\n", timestamp, level, category, message, errStr)
	return []byte(logLine), nil
}

var consoleLogger = logrus.New()

// initConsoleLogger configures consoleLogger to match the CARBON_CONSOLE settings.
func initConsoleLogger() {
	consoleLogger.SetOutput(os.Stdout)                // Log to console.
	consoleLogger.SetLevel(logrus.DebugLevel)           // Threshold set to DEBUG.
	consoleLogger.SetFormatter(&CustomFormatter{})      // Use our custom formatter.
}




// # CARBON_CONSOLE is set to be a ConsoleAppender using a PatternLayout.
// appender.CARBON_CONSOLE.type = Console
// appender.CARBON_CONSOLE.name = CARBON_CONSOLE
// appender.CARBON_CONSOLE.layout.type = PatternLayout
// appender.CARBON_CONSOLE.layout.pattern = [%d] %5p {%c{1}} - %m%ex%n
// appender.CARBON_CONSOLE.filter.threshold.type = ThresholdFilter
// appender.CARBON_CONSOLE.filter.threshold.level = DEBUG


// Zap code for the same:

// package main

// import (
//     "fmt"
//     "os"
//     "strings"
//     "time"

//     "go.uber.org/zap"
//     "go.uber.org/zap/zapcore"
// )

// // CustomTimeEncoder formats the time as [YYYY-MM-DD HH:MM:SS]
// func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//     enc.AppendString(fmt.Sprintf("[%s]", t.Format("2006-01-02 15:04:05")))
// }

// // PaddedLevelEncoder pads the level to 5 characters
// func PaddedLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
//     enc.AppendString(fmt.Sprintf("%-5s", level.CapitalString()))
// }

// // ShortCallerEncoder logs the last segment of the caller
// func ShortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
//     parts := strings.Split(caller.TrimmedPath(), "/")
//     enc.AppendString(fmt.Sprintf("{%s}", parts[len(parts)-1]))
// }

// func main() {
//     encoderConfig := zapcore.EncoderConfig{
//         TimeKey:        "time",
//         LevelKey:       "level",
//         NameKey:        "logger",
//         CallerKey:      "caller",
//         MessageKey:     "msg",
//         StacktraceKey:  "stacktrace",
//         LineEnding:     zapcore.DefaultLineEnding,
//         EncodeTime:     CustomTimeEncoder,
//         EncodeLevel:    PaddedLevelEncoder,
//         EncodeCaller:   ShortCallerEncoder,
//         EncodeName:     zapcore.FullNameEncoder,
//     }

//     core := zapcore.NewCore(
//         zapcore.NewConsoleEncoder(encoderConfig),
//         zapcore.AddSync(os.Stdout),
//         zap.DebugLevel,
//     )

//     logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

//     logger.Debug("This is a debug message.")
//     logger.Info("This is an info message.")
//     logger.Error("This is an error message.")
// }
