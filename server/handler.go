package server

import (
	"log/slog"

	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/cmd"
)

func (s *SSHServer) handleCmd(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		cli := cmd.NewCLI(s.cfg, sess)
		cli.Init()
		if err := cli.Run(); err != nil {
			slog.Error("Could not run command", "error", err)
		}

		next(sess)
	}
}
