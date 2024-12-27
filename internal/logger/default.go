package logger

import (
	"log/slog"
	"os"
)

func SetDefault() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(log)
}
