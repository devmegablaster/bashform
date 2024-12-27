package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/services"
)

type SSHServer struct {
	cfg       config.Config
	db        *database.Database
	userSvc   *services.UserService
	cryptoSvc *services.CryptoService
	s         *ssh.Server
	wg        *sync.WaitGroup
}

func (s *SSHServer) initSSHServer() error {
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(s.cfg.SSH.Host, strconv.Itoa(s.cfg.SSH.Port))),
		wish.WithHostKeyPath(s.cfg.SSH.KeyPath),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return s.handleAuth(ctx, key)
		}),
		wish.WithMiddleware(
			s.handleCmd,
			activeterm.Middleware(),
			s.Logger,
		),
	)

	if err != nil {
		return err
	}

	s.s = server
	return nil
}

func NewSSHServer(wg *sync.WaitGroup, cfg config.Config, db *database.Database) (*SSHServer, error) {
	s := &SSHServer{
		cfg:       cfg,
		db:        db,
		userSvc:   services.NewUserService(cfg, db),
		cryptoSvc: services.NewCryptoService(cfg.Crypto),
		wg:        wg,
	}

	if err := s.initSSHServer(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SSHServer) ListenAndServe(ctx context.Context) error {
	errChan := make(chan error, 1)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		slog.Info("✅ Starting SSH server", slog.String("address", s.s.Addr))
		if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		slog.Error("❌ Could not start server", "error", err)
		return err
	case <-ctx.Done():
		slog.Info("🔌 Shutting down SSH Server")
		return s.shutdown()
	}
}

func (s *SSHServer) shutdown() error {
	timeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	s.wg.Add(1)
	defer s.wg.Done()

	// wait for all connections to close or timeout
	return s.s.Shutdown(timeCtx)
}
