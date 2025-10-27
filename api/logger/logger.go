package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
    var handler slog.Handler
    
    if env == "production" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        })
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
            AddSource: true, // file:line
        })
    }
    
    return slog.New(handler)
}