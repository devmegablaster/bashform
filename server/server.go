package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/devmegablaster/bashform/internal/config"
)

type SSHServer struct {
	cfg config.Config
	s   *ssh.Server
}

func (s *SSHServer) initSSHServer() {
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(s.cfg.SSH.Host, strconv.Itoa(s.cfg.SSH.Port))),
		wish.WithHostKeyPath(s.cfg.SSH.KeyPath),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true // TODO: Implement auth here
		}),
		wish.WithMiddleware(
			s.handleCmd,
			activeterm.Middleware(),
			s.Logger,
		),
	)

	if err != nil {
		slog.Error("Could not start server", "error", err)
	}

	s.s = server
}

func NewSSHServer(cfg config.Config) *SSHServer {
	s := &SSHServer{
		cfg: cfg,
	}
	s.initSSHServer()

	return s
}

func (s *SSHServer) ListenAndServe(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		slog.Info("Starting SSH server", slog.String("address", s.s.Addr))
		if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		slog.Error("Could not start server", "error", err)
		return err
	case <-ctx.Done():
		timeCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.s.Shutdown(timeCtx)
	}
}
