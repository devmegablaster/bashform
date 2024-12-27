package services

import (
	"fmt"
	"log/slog"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
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

// create a new user using the user request
func (s *UserService) Create(userReq models.UserRequest) (*models.User, error) {
	user := userReq.ToUser()
	err := s.ur.Create(&user)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, fmt.Errorf("Failed to create user")
	}

	return &user, nil
}

// get a user by their public key
func (s *UserService) GetByPubKey(pubKey string) (*models.User, error) {
	user, err := s.ur.GetByPubKey(pubKey)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return nil, fmt.Errorf("Failed to get user")
	}

	return user, nil
}

// update a user using user strcut
func (s *UserService) Update(user *models.User) error {
	err := s.ur.Update(user)
	if err != nil {
		slog.Error("Failed to update user", "error", err)
		return fmt.Errorf("Failed to update user")
	}

	return nil
}
