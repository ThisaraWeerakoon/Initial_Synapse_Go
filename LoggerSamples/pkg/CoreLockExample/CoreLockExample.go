package main

import (
        "fmt"
        "os"
        "sync"
        "go.uber.org/zap"
        "go.uber.org/zap/zapcore"
)

// Package-level loggers
var (
		mainLog *zap.Logger
        packageALog *zap.Logger
        packageBLog *zap.Logger
        packageCLog *zap.Logger
)

var config zap.Config

func init() {
	config = zap.NewProductionConfig()
	logFile, err:= os.OpenFile("./combined.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644) // Open the file
	if err!= nil {
			panic(err) // Handle error appropriately
	}

	// Use zapcore.NewSyncWriter to write to the file.
	writer:= zapcore.AddSync(logFile)

	mainCore:= zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			writer, // Write to the file
			zap.InfoLevel,
	)
	mainLog = zap.New(mainCore)

	packageACore:= zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			writer, // Write to the file
			zap.DebugLevel,
	)
	packageALog = zap.New(packageACore)

	packageBCore:= zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			writer, // Write to the file
			zap.WarnLevel,
	)
	packageBLog = zap.New(packageBCore)

	packageCCore:= zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			writer, // Write to the file
			zap.ErrorLevel,
	)
	packageCLog = zap.New(packageCCore)
}

func main() {
        defer mainLog.Sync()
        defer packageALog.Sync()
        defer packageBLog.Sync()
        defer packageCLog.Sync()

        longMessage:= generateLongMessage(10000) // 10,000 characters

        var wg sync.WaitGroup
        numRoutines:= 100 // Number of concurrent writers

        for i:= 0; i < numRoutines; i++ {
                wg.Add(1)
                go func(id int) {
                        defer wg.Done()
                        packageALog.Debug(fmt.Sprintf("Package A: %s (Routine %d)", longMessage, id))
                        packageBLog.Warn(fmt.Sprintf("Package B: %s (Routine %d)", longMessage, id))
                        packageCLog.Error(fmt.Sprintf("Package C: %s (Routine %d)", longMessage, id))
                }(i)
        }

        wg.Wait()
    fmt.Println("Writing is finished")

}

func generateLongMessage(length int) string {
        message := make([]byte, length)
        for i:= 0; i < length; i++ {
                message[i] = byte('A' + (i % 26)) // Use A-Z characters
        }
        return string(message)
}