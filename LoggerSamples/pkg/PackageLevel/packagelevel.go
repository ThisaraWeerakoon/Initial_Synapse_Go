package packagelevel


import (
        "go.uber.org/zap"
        "go.uber.org/zap/zapcore"
        "os"
)

// Package-level loggers (one per package/component)
var (
        dataNucleusLog *zap.Logger
        hiveLog        *zap.Logger
        mainLog        *zap.Logger // For the main application
)

// Global configuration (can be customized)
var config zap.Config

// Initialize loggers (ideally in an init() function)
func init() {
        // Configuration for the main logger (can be customized)
        config = zap.NewProductionConfig() // Or zap.NewDevelopmentConfig() for dev
    config.OutputPaths = []string{"stdout", "my_app.log"} // Log to console and file


        // Create the main logger
        mainLog, _ = config.Build()

        // Create separate cores and loggers for DataNucleus and Hive
        dataNucleusCore:= zapcore.NewCore(
                zapcore.NewJSONEncoder(config.EncoderConfig), // Use the same encoder config
                zapcore.Lock(os.Stdout),                     // Share the same output (or different)
                zap.ErrorLevel,                               // DataNucleus level: ERROR
        )
        dataNucleusLog = zap.New(dataNucleusCore)

        hiveCore:= zapcore.NewCore(
                zapcore.NewJSONEncoder(config.EncoderConfig), // Use the same encoder config
                zapcore.Lock(os.Stdout),                     // Share the same output (or different)
                zap.WarnLevel,                               // Hive level: WARN
        )
        hiveLog = zap.New(hiveCore)

}

func main() {
        defer mainLog.Sync()   // Flush main logs at the end
        defer dataNucleusLog.Sync() // Flush DataNucleus logs
        defer hiveLog.Sync()   // Flush Hive logs

        mainLog.Info("Application started")

        dataNucleusLog.Error("DataNucleus error occurred") // This WILL be logged

        hiveLog.Info("Hive info message")   // This will NOT be logged (level is WARN)
        hiveLog.Warn("Hive warning message") // This WILL be logged

        mainLog.Debug("Main app debug message") // This may or may not be logged depending on the global level.

    // Example of using fields for more structured logging:
    mainLog.With(zap.String("component", "main")).Info("Main app info with fields")
    dataNucleusLog.With(zap.String("component", "DataNucleus")).Error("DataNucleus error with fields")
    hiveLog.With(zap.String("component", "Hive")).Warn("Hive warning with fields")

}

// Example of a helper function for logging with a specific level:
func DataNucleusLog(level zapcore.Level, msg string, fields...zap.Field) {
    if dataNucleusLog.Core().Enabled(level) {
        dataNucleusLog.Log(level, msg, fields...)
    }
}

// Example usage of the helper function:
func someDataNucleusFunction() {
    DataNucleusLog(zap.DebugLevel, "This is a debug message from DataNucleus", zap.String("key", "value")) // Won't be logged
    DataNucleusLog(zap.ErrorLevel, "This is a error message from DataNucleus", zap.String("key", "value"))  // Will be logged
}