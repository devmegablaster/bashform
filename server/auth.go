package server

import (
	"log/slog"

	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/models"
)

func (s *SSHServer) handleAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	encodedKey := s.cryptoSvc.Base64Encode(key.Marshal())
	user, err := s.userSvc.GetByPubKey(encodedKey)
	if err != nil {
		userReq := models.UserRequest{
			PubKey: encodedKey,
		}
		user, err = s.userSvc.Create(userReq)
		if err != nil {
			return false
		}

		slog.Info("🔑 Created new user", slog.String("user", user.ID.String()))
	}

	ctx.SetValue("user", user)
	return true
}
