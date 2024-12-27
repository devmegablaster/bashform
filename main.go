package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/server"
)

func main() {

	cfg := config.New()

	s := server.NewSSHServer(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		if err := s.ListenAndServe(ctx); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		slog.Error("Could not start server", "error", err)
	case sig := <-sigChan:
		slog.Info("Shutting down", slog.String("signal", sig.String()))
		cancel()
	}
}
