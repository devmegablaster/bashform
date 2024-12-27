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

func (s *SSHServer) addUsername(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		user := sess.Context().Value("user").(*models.User)
		if user.Name == "" && sess.User() != "" {
			user.Name = sess.User()
			if err := s.userSvc.Update(user); err != nil {
				slog.Error("Failed to update user", "error", err)
				return
			}

			slog.Info("🔑 Updated user name", slog.String("user", user.ID.String()), slog.String("name", user.Name))
		}

		next(sess)
	}
}
