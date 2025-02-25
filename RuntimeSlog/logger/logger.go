package logger

import (
        "log/slog"
        "os"

        "github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

        "github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/models"

)

// Define an interface for setting log levels.
type LogLevelSetter interface {
        SetLevel(k *koanf.Koanf, name string)
}

// Global koanf instance.
var k = koanf.New(".")

var currentK = koanf.New(".")

var koanfChan = make(chan *koanf.Koanf)

var loggers = make(map[string]LogLevelSetter) // Store LogLevelSetters

func init() {

        go notifier(koanfChan)
        // file reader
        f := file.Provider("conf.json")

        // Load initial configuration. For simplicity we only takes the log level from the .json file. Future we extend this into the entire configuration.
        if err := k.Load(f, json.Parser()); err != nil {
                slog.Error("error loading initial config", "error", err)
                os.Exit(1) // Exit if initial config load fails
        }

        currentK = k

        // Watch the file and get a callback on change. The callback re-load the
        // configuration.
        f.Watch(func(event interface{}, err error) {
                if err != nil {
                        slog.Error("error watching config file", "error", err)
                        return
                }

                newK := koanf.New(".") // Create a new Koanf instance
                if err := newK.Load(f, json.Parser()); err != nil {
                        slog.Error("error reloading config", "error", err)
                        return
                }

                k = newK // Atomically swap the Koanf instance

                //callback function
                koanfChan <- k

        })

}

// func RegisterLogger(name string, setter LogLevelSetter) {
//         loggers[name] = setter
// }
    

func InitializeLogger(name string) *models.CustomLogger {
        customLogger := models.CustomLogger{}
        customLogger.SetLevel(k, name) // Initial level set

        loggers[name] = &customLogger // Store the logger
        return &customLogger
}

func notifier(dataChan <-chan *koanf.Koanf) {
        for newK := range dataChan {
                for name, setter := range loggers {
                        if newK.String(name) != currentK.String(name) {
                                setter.SetLevel(newK, name)
                        }
                }
                currentK = newK // Update currentK after processing all loggers
        }
}