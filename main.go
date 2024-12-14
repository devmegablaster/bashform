package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/devmegablaster/bashform/internal/config"
)

func main() {

	cfg := config.New()

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.SSH.Host, strconv.Itoa(cfg.SSH.Port))),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					// TODO: Application Here

					next(sess)
				}
			},
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", cfg.SSH.Host, "port", cfg.SSH.Port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
