package main

import (
        "log/slog"
        "os"
        "github.com/ThisaraWeerakoon/Initial_Synapse_Go/PackageLevelSlog/mypackage" // Import your package
)


func main() {
    // Global logger (can be minimal)
    globalLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo, // Default level for other packages if needed.
    }))

    globalLogger.Info("Application started") // Example

    mypackage.MyFunction() // Call function in the other package

}
