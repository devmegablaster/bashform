package server

import (
	"github.com/charmbracelet/ssh"
)

func (s *SSHServer) handleAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	encodedKey := s.cryptoSvc.Base64Encode(key.Marshal())
	user, err := s.userSvc.GetByPubKey(encodedKey)
	if err != nil {
		return false
	}

	ctx.SetValue("user", user)
	return true
}
