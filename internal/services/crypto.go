package services

import (
	"bytes"
	"encoding/base64"

	"github.com/devmegablaster/bashform/internal/config"
)

type CryptoService struct {
	cfg config.CryptoConfig
}

func NewCryptoService(cfg config.CryptoConfig) *CryptoService {
	return &CryptoService{
		cfg: cfg,
	}
}

func (c *CryptoService) Base64Encode(data []byte) string {
	encoded := bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, &encoded)
	enc.Write(data)

	return encoded.String()
}
