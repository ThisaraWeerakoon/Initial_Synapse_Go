package carbonlogfilewithoutrotatelogs

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
)

func CarbonLogfileWithoutRotatelogs() {
	// Configure lumberjack for file-based logging:
	// - Filename: "logs/wso2carbon.log" (active log file).
	// - MaxSize: 10 MB per file.
	// - MaxBackups: keep up to 20 backup files.
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/wso2carbon.log",
		MaxSize:    10,  // megabytes
		MaxBackups: 20,  // maximum number of backup files
		// Note: lumberjack does not support time-based rotation out of the box.
	}

	// Create zap WriteSyncers: one for file (via lumberjack) and one for console.
	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)
	multiWriter := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)

	// Configure zap's encoder to mimic our desired log format.
	// The pattern roughly corresponds to:
	// "[timestamp] LEVEL {logger} - message"
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		MessageKey: "msg",
		CallerKey:  "caller",
		// Use a custom time layout similar to Log4j's timestamp.
		EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		// CapitalLevelEncoder produces levels like DEBUG, INFO, etc.
		EncodeLevel: zapcore.CapitalLevelEncoder,
		// Caller info can be trimmed to show only the file name and line number.
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// Use a console encoder for human-readable output.
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Set the minimum logging level to DEBUG.
	logLevel := zap.DebugLevel

	// Create the zap core.
	core := zapcore.NewCore(encoder, multiWriter, logLevel)

	// Build the logger with caller and stacktrace options.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	// Start a goroutine to trigger daily rotation at midnight.
	go scheduleDailyRotation(lumberjackLogger)

	// Log sample messages in a loop.
	for {
		logger.Debug("Application is running...", zap.String("logger", "main"))
		time.Sleep(10 * time.Second)
	}
}

// scheduleDailyRotation calculates the time until the next midnight,
// then sleeps until that time before triggering lumberjack's Rotate method.
func scheduleDailyRotation(lj *lumberjack.Logger) {
	for {
		now := time.Now()
		// Compute the next midnight (00:00:00).
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(nextMidnight))
		lj.Rotate()
	}
}
