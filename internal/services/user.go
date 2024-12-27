package services

import (
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
	"github.com/devmegablaster/bashform/internal/types"
)

type UserService struct {
	cfg config.Config
	db  *database.Database
	ur  *repository.UserRepository
}

func NewUserService(cfg config.Config, db *database.Database) *UserService {
	return &UserService{
		cfg: cfg,
		db:  db,
		ur:  repository.NewUserRepository(db),
	}
}

func (s *UserService) Create(user *models.User) error {
	err := s.ur.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByPubKey(pubKey string) (*models.User, types.ServiceErrors) {
	user, err := s.ur.GetByPubKey(pubKey)
	if err != nil {
		return nil, types.ServiceErrors{"error": err.Error()}
	}

	return user, nil
}
