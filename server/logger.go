package server

import (
	"log/slog"
	"time"

	"github.com/charmbracelet/ssh"
)

func (s *SSHServer) logger(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		init := time.Now()
		slog.Info("New Connection", "remote_addr", sess.RemoteAddr(), "username", sess.User(), "command", sess.Command())
		next(sess)
		duration := time.Since(init)
		slog.Info("Connection Closed", "remote_addr", sess.RemoteAddr(), "username", sess.User(), "duration", duration.Round(time.Second).String())
	}
}
