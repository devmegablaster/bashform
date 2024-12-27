package server

import (
	"log/slog"

	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/cmd"
)

func (s *SSHServer) handleCmd(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		if err := s.executeCommand(sess); err != nil {
			slog.Error("Command execution failed", "error", err)
			return
		}
		next(sess)
	}
}

func (s *SSHServer) executeCommand(sess ssh.Session) error {
	cli := cmd.NewCLI(s.cfg, s.db, sess)
	return cli.Run()
}
