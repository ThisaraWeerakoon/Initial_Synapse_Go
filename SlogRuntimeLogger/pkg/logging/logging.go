package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type Logger struct{
	PackageName string
	LevelMap *map[string]string
	Handler slog.Handler //stoudt,json type 

	// DefaultLevel slog.Leveler
}

// Erros should be defined and other handler building logic should be revised 
func NewLogger(levelMap *map[string]string, handlerConfig HandlerConfig) (Logger,error){
	l := Logger{}
	format := handlerConfig.Format
	outputPath := handlerConfig.OutputPath
	switch format {
	case "json":
		switch outputPath {
		case "stdout":
			l.Handler = slog.NewJSONHandler(os.Stdout, nil)
		case "file":
			// l.Handler = slog.NewJSONHandler(slog.File(outputPath), slog.DefaultTimeFormat)
		}
	case "text":
		switch outputPath {
		case "stdout":
			l.Handler = slog.NewTextHandler(os.Stdout, nil)
		case "file":
			// l.Handler = slog.NewTextHandler(slog.File(outputPath), slog.DefaultTimeFormat)
		}
	}

	l.LevelMap = levelMap
	return l,nil

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

//Type to hold the handler related configurations from Config 
type HandlerConfig struct {
	//json,text
	Format string `koanf:"format"`
	//stdout, file
	OutputPath string `koanf:"outputPath"`
	

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



func (l *Logger) Debugw(msg string, args ...any){
	if l.LevelMap == nil {
		//Handle the case where the LevelMap is nil. Either return, or log an error and use a default.
		slog.Error("LevelMap is nil in Debugw")
		logger := slog.New(NewLevelHandler(slog.LevelDebug, l.Handler)) // or some default.
		logger.Debug(msg, args...)
		return
	}

	levelStr, ok := (*l.LevelMap)[l.PackageName] // Dereference the pointer and access the map.
	if !ok {
			// Handle the case where the PackageName is not in the map.
			slog.Error("PackageName not found in LevelMap", "PackageName", l.PackageName)
			logger := slog.New(NewLevelHandler(slog.LevelDebug, l.Handler)) //or some default.
			logger.Debug(msg, args...)
			return
	}

	logger := slog.New(NewLevelHandler(LevelFromString(levelStr), l.Handler))
	logger.Debug(msg, args...)
	
}

func (l *Logger) Infow(msg string, args ...any){
	if l.LevelMap == nil {
		//Handle the case where the LevelMap is nil. Either return, or log an error and use a default.
		slog.Error("LevelMap is nil in Infow")
		logger := slog.New(NewLevelHandler(slog.LevelInfo, l.Handler)) // or some default.
		logger.Info(msg, args...)
		return
	}

	levelStr, ok := (*l.LevelMap)[l.PackageName] // Dereference the pointer and access the map.
	if !ok {
			// Handle the case where the PackageName is not in the map.
			slog.Error("PackageName not found in LevelMap", "PackageName", l.PackageName)
			logger := slog.New(NewLevelHandler(slog.LevelInfo, l.Handler)) //or some default.
			logger.Info(msg, args...)
			return
	}

	logger := slog.New(NewLevelHandler(LevelFromString(levelStr), l.Handler))
	logger.Info(msg, args...)
	
}

func (l *Logger) Warnw(msg string, args ...any){
	if l.LevelMap == nil {
		//Handle the case where the LevelMap is nil. Either return, or log an error and use a default.
		slog.Error("LevelMap is nil in Warnw")
		logger := slog.New(NewLevelHandler(slog.LevelWarn, l.Handler)) // or some default.
		logger.Warn(msg, args...)
		return
	}

	levelStr, ok := (*l.LevelMap)[l.PackageName] // Dereference the pointer and access the map.
	if !ok {
			// Handle the case where the PackageName is not in the map.
			slog.Error("PackageName not found in LevelMap", "PackageName", l.PackageName)
			logger := slog.New(NewLevelHandler(slog.LevelWarn, l.Handler)) //or some default.
			logger.Warn(msg, args...)
			return
	}

	logger := slog.New(NewLevelHandler(LevelFromString(levelStr), l.Handler))
	logger.Warn(msg, args...)
	
}

func (l *Logger) Errorw(msg string, args ...any){
	if l.LevelMap == nil {
		//Handle the case where the LevelMap is nil. Either return, or log an error and use a default.
		slog.Error("LevelMap is nil in Errorw")
		logger := slog.New(NewLevelHandler(slog.LevelError, l.Handler)) // or some default.
		logger.Error(msg, args...)
		return
	}

	levelStr, ok := (*l.LevelMap)[l.PackageName] // Dereference the pointer and access the map.
	if !ok {
			// Handle the case where the PackageName is not in the map.
			slog.Error("PackageName not found in LevelMap", "PackageName", l.PackageName)
			logger := slog.New(NewLevelHandler(slog.LevelError, l.Handler)) //or some default.
			logger.Error(msg, args...)
			return
	}

	logger := slog.New(NewLevelHandler(LevelFromString(levelStr), l.Handler))
	logger.Error(msg, args...)
	
}