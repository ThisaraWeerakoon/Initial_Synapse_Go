package logger

import (
        "log/slog"
        "os"

        "github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

        "github.com/ThisaraWeerakoon/Initial_Synapse_Go/RuntimeSlog/models"
)

// Global koanf instance.
var k = koanf.New(".")

var currentK = koanf.New(".")

var koanfChan = make(chan *koanf.Koanf)

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

func InitializeLogger(name string) *models.CustomLogger {
        packageALogger := models.CustomLogger{}
        packageALogger.SetLevel(k, name)
        return &packageALogger
}