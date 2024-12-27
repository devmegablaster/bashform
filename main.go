package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/logger"
	"github.com/devmegablaster/bashform/server"
)

func main() {

	logger.SetDefault()

	cfg := config.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	db, err := database.New(ctx, wg, cfg.Database)
	if err != nil {
		slog.Error("❌ Could not connect to database", "error", err)
		os.Exit(1)
	}

	s, err := server.NewSSHServer(wg, cfg, db)
	if err != nil {
		slog.Error("❌ Could not create server", "error", err)
		os.Exit(1)
	}

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
		slog.Error("❌ Could not start server", "error", err)
		os.Exit(1)

	case sig := <-sigChan:
		slog.Info("🚨 Shutting down gracefully", slog.String("signal", sig.String()))
		cancel()
		wg.Wait()
		slog.Info("✅ Graceful shutdown complete")
	}
}
