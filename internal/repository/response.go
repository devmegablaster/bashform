package repository

import (
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
)

type ResponseRepository struct {
	db *database.Database
}

func NewResponseRepository(db *database.Database) *ResponseRepository {
	return &ResponseRepository{
		db: db,
	}
}

func (r *ResponseRepository) Create(response *models.Response) error {
	if err := r.db.DB.Create(&response).Error; err != nil {
		return err
	}

	return nil
}

func (r *ResponseRepository) GetByUserAndFormID(userID, formID string) (*models.Response, error) {
	var response models.Response

	if err := r.db.DB.Preload("Answers").First(&response, "user_id = ? AND form_id = ?", userID, formID).Error; err != nil {
		return nil, err
	}

	return &response, nil
}
