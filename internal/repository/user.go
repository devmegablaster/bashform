package repository

import (
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByPubKey(pubKey string) (*models.User, error) {
	var user models.User

	if err := r.db.DB.First(&user, "pub_key = ?", pubKey).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
