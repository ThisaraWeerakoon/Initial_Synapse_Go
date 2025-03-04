package main

import (
	"SlogRuntimeLogger/pkg/config"
	"SlogRuntimeLogger/pkg/logging"
	"SlogRuntimeLogger/pkg/packageA"
	"SlogRuntimeLogger/pkg/packageB"
	"context"
	"fmt"
	"log"
	"time"
)


func main() {
	configFilePath := "config.yaml"
	cfg, err := config.ReadFile(configFilePath)
 	if err != nil {
			log.Fatalf("Cannot read config file: %s", err.Error())
	}


	var levelMap *map[string]string
	var handlerConfig logging.HandlerConfig
	if cfg.IsSet("logger") {
		handlerConfig = logging.HandlerConfig{}
		//Initialize handlerConfig
		cfg.MustUnmarshal("logger.handler", &handlerConfig)
		//Initialize levelMap
		levelMap = &map[string]string{} //This level map is used by entire application
		cfg.MustUnmarshal("logger.level.packages", levelMap)
	}

	logger, err := logging.NewLogger(levelMap, handlerConfig)

	if err != nil {
		log.Fatalf("Cannot create logger: %s", err.Error())
	}

	packageA := packageA.New("A", "main", logger)
	packageB := packageB.New("B", "main", logger)	

	cfg.Watch(context.Background(),configFilePath,levelMap)

	for {
		fmt.Println("Looping............................................................")
		packageA.Run()
		packageB.Run()
		time.Sleep(5 * time.Second)
	}

	


    
	// logger, err := logging.NewNamedFromConfig(loggerConfig, "main")
	// if err != nil {
	// 	log.Fatalf("Cannot create logger: %s", err.Error())
	// }

	// packageA := packageA.New("John", 30, logger)
	// packageB := packageB.New("Jane", 25, logger)

	// cfg.Watch(context.Background(),configFilePath,loggerConfig)
	// for {
	// 	fmt.Println("Logging level at main is:", logger.Level())
	// 	logger.Debug("This is a debug message")
	// 	logger.Info("This is an info message")
	// 	// logger.Warn("This is a warning message")
	// 	// logger.Error("This is an error message")
	// 	packageA.Run()
	// 	packageB.Run()
	// 	time.Sleep(5 * time.Second)
	// }



	// Block until a signal is received.
	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// sig := <-sigCh
	// log.Printf("Received signal: %v. Exiting...", sig)
}













// package main

// import (
// 	"context"
// 	"log/slog"
// 	"os"
// )

// // A LevelHandler wraps a Handler with an Enabled method
// // that returns false for levels below a minimum.
// type LevelHandler struct {
// 	level   slog.Leveler
// 	handler slog.Handler
// }

// // NewLevelHandler returns a LevelHandler with the given level.
// // All methods except Enabled delegate to h.
// func NewLevelHandler(level slog.Leveler, h slog.Handler) *LevelHandler {
// 	// Optimization: avoid chains of LevelHandlers.
// 	if lh, ok := h.(*LevelHandler); ok {
// 		h = lh.Handler()
// 	}
// 	return &LevelHandler{level, h}
// }

// // Enabled implements Handler.Enabled by reporting whether
// // level is at least as large as h's level.
// func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
// 	return level >= h.level.Level()
// }

// // Handle implements Handler.Handle.
// func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
// 	return h.handler.Handle(ctx, r)
// }

// // WithAttrs implements Handler.WithAttrs.
// func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
// 	return NewLevelHandler(h.level, h.handler.WithAttrs(attrs))
// }

// // WithGroup implements Handler.WithGroup.
// func (h *LevelHandler) WithGroup(name string) slog.Handler {
// 	return NewLevelHandler(h.level, h.handler.WithGroup(name))
// }

// // Handler returns the Handler wrapped by h.
// func (h *LevelHandler) Handler() slog.Handler {
// 	return h.handler
// }

// // This example shows how to Use a LevelHandler to change the level of an
// // existing Handler while preserving its other behavior.
// //
// // This example demonstrates increasing the log level to reduce a logger's
// // output.
// //
// // Another typical use would be to decrease the log level (to LevelDebug, say)
// // during a part of the program that was suspected of containing a bug.
// func main() {
// 	th := slog.NewTextHandler(os.Stdout, nil)
// 	logger := slog.New(NewLevelHandler(slog.LevelDebug, th))
// 	logger.Debug("printed")
// 	logger.Info("not printed")
// 	logger.Warn("printed")
// 	logger.Error("printed")

// }