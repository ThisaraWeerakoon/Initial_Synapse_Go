package carbonlogfile

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// initCarbonLogfileLogger creates a zap logger that writes to a daily rotated file.
// It mimics the log4j configuration:
//   - Log file: /path/to/wso2carbon.log (current log file via symlink)
//   - File pattern: wso2carbon-%m-%d-%Y.log (daily rotation)
//   - Rotation time: 24 hours (daily)
//   - Max backups: 20 files
//   - Threshold: DEBUG and above
func initCarbonLogfileLogger() (*zap.Logger, error) {
	// Path to the active log file and the pattern for rotated files.
	logPath := "/path/to/wso2carbon.log" // Change this to your logfiles.home path.
	pattern := "/path/to/wso2carbon-%m-%d-%Y.log"

	// Set up rotatelogs: daily rotation with a maximum of 20 rotated files.
	rotator, err := rotatelogs.New(
		pattern,
		rotatelogs.WithLinkName(logPath),        // create symlink to current log file
		rotatelogs.WithRotationTime(24*time.Hour), // rotate every day
		rotatelogs.WithRotationCount(20),          // keep max 20 old files
	)
	if err != nil {
		return nil, err
	}

	// Create a zap encoder configuration to mimic the PatternLayout "[%d] %5p {%c} - %m%ex%n"
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger", // this will be the category (%c)
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // e.g. DEBUG, INFO, etc.
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
	}

	// Use a console encoder (human-readable) for simplicity.
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create a zap core that writes to the rotator with a minimum level of DEBUG.
	core := zapcore.NewCore(encoder, zapcore.AddSync(rotator), zap.DebugLevel)

	// Build and return the logger.
	logger := zap.New(core)
	return logger, nil
}


