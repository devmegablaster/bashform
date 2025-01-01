package repository

import (
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
)

type FormRepository struct {
	db *database.Database
}

func NewFormRepository(db *database.Database) *FormRepository {
	return &FormRepository{
		db: db,
	}
}

func (r *FormRepository) Create(form *models.Form) error {
	if err := r.db.DB.Create(&form).Error; err != nil {
		return err
	}

	return nil
}

func (r *FormRepository) GetForUser(userID string) ([]models.Form, error) {
	var forms []models.Form

	if err := r.db.DB.Find(&forms, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return forms, nil
}

func (r *FormRepository) GetByCode(code string) (*models.Form, error) {
	var form models.Form

	if err := r.db.DB.Preload("Questions.Options").First(&form, "code = ?", code).Error; err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) GetByID(id string) (*models.Form, error) {
	var form models.Form

	if err := r.db.DB.Preload("Questions").First(&form, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) GetWithResponses(userID, formID string) (*models.Form, error) {
	var form models.Form

	if err := r.db.DB.Preload("Responses.Answers").Preload("Questions").First(&form, "id = ? AND user_id = ?", formID, userID).Error; err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) GetWithResponsesUsingCode(userID, code string) (*models.Form, error) {
	var form models.Form

	if err := r.db.DB.Preload("Responses.Answers").Preload("Questions").First(&form, "code = ? AND user_id = ?", code, userID).Error; err != nil {
		return nil, err
	}

	return &form, nil
}
