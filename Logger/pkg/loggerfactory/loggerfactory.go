package loggerfactory

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

//Type to extract and hold the slog handler related configurations from Config
// format : json/text
// outputpath: stdout/file/stderr
type SlogHandlerConfig struct {
	//json,text
	Format string `koanf:"format"`
	//stdout, file
	OutputPath string `koanf:"outputPath"`
}

//Intentionally put 'slog' in future we can introduce more abstract handlers. Every handler should implement slog.Handler interface
func GetSlogHandler(slogHandlerConfig SlogHandlerConfig) slog.Handler {
	format := slogHandlerConfig.Format
	outputPath := slogHandlerConfig.OutputPath
	var slogHandler slog.Handler
	switch format {
	case "json":
		switch outputPath {
		case "stdout":
			slogHandler = slog.NewJSONHandler(os.Stdout, nil)
		case "file":
			// l.Handler = slog.NewJSONHandler(slog.File(outputPath), slog.DefaultTimeFormat)
		}
	case "text":
		switch outputPath {
		case "stdout":
			slogHandler = slog.NewTextHandler(os.Stdout, nil)
		case "file":
			// l.Handler = slog.NewTextHandler(slog.File(outputPath), slog.DefaultTimeFormat)
		}
	}
	return slogHandler
} 


// A LevelHandler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type LevelHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

// NewLevelHandler returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func NewLevelHandler(level slog.Leveler, h slog.Handler) *LevelHandler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*LevelHandler); ok {
		h = lh.Handler()
	}
	return &LevelHandler{level, h}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements Handler.Handle.
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithGroup(name))
}

// Handler returns the Handler wrapped by h.
func (h *LevelHandler) Handler() slog.Handler {
	return h.handler
}

// LevelFromString converts a string representation of a log level to a slog.Leveler.
func LevelFromString(levelStr string) slog.Leveler {
	switch strings.ToLower(levelStr) {
	case "debug":
			return slog.LevelDebug
	case "info":
			return slog.LevelInfo
	case "warn", "warning":
			return slog.LevelWarn
	case "error":
			return slog.LevelError
	default:
			// Return default level (e.g., Info) or handle invalid input as needed.
			return slog.LevelInfo // Or return an error, or a custom level.
	}
}



func GetLogger(packageName string, logLevelMap *map[string]string, slogHandlerConfig SlogHandlerConfig) *slog.Logger{
	levelStr, ok := (*logLevelMap)[packageName] // Dereference the pointer and access the map.
	if !ok {
			// Handle the case where the PackageName is not in the map.
			slog.Error("PackageName not found in LevelMap", "PackageName", packageName)
			logger := slog.New(NewLevelHandler(slog.LevelDebug,GetSlogHandler(slogHandlerConfig))) //or some default.
			return logger
	}
	return slog.New(NewLevelHandler(LevelFromString(levelStr), GetSlogHandler(slogHandlerConfig)))

}

