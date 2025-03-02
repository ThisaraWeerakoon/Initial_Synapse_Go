package main

import (
	"log/slog"
	"time"

	"github.com/natefinch/lumberjack"
)

// limitedWriter wraps a lumberjack.Logger and counts bytes.
// When the written bytes exceed 1KB, it triggers a rotation.
type limitedWriter struct {
	lw    *lumberjack.Logger
	count int
}

func (w *limitedWriter) Write(p []byte) (n int, err error) {
	n, err = w.lw.Write(p)
	w.count += n
	if w.count >= 1024 { // 1KB threshold reached
		// Force rotation.
		_ = w.lw.Rotate()
		// Reset counter after rotation.
		w.count = 0
	}
	return n, err
}

func main() {
	// Set up lumberjack (MaxSize is still 1 MB, but our wrapper will trigger at 1KB).
	lj := &lumberjack.Logger{
		Filename:   "test.log",
		MaxSize:    1, // 1 MB (but will not be reached due to our manual rotation)
		MaxBackups: 3,
		MaxAge:     7, // days
	}

	// Wrap lumberjack with our limitedWriter to force rotation at ~1KB.
	lw := &limitedWriter{lw: lj}

	// Create a slog logger that writes JSON-formatted logs to our limitedWriter.
	logger := slog.New(slog.NewJSONHandler(lw, nil))

	// Write many log messages to trigger rotation.
	for i := 0; i < 10000; i++ {
		logger.Info("Test log message", "iteration", i)
		time.Sleep(10 * time.Millisecond)
	}
}
