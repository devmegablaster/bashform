package server

import (
	"log/slog"

	"github.com/charmbracelet/ssh"
)

func (s *SSHServer) Logger(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		slog.Info("New Connection", "remote_addr", sess.RemoteAddr(), "user", sess.User(), "command", sess.Command())
		next(sess)
		slog.Info("Connection Closed", "remote_addr", sess.RemoteAddr(), "user", sess.User())
	}
}
