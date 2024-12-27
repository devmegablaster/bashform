package services

import (
	"errors"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
	"github.com/devmegablaster/bashform/internal/types"
	"gorm.io/gorm"
)

type FormService struct {
	cfg config.Config
	db  *database.Database
	fr  *repository.FormRepository
	rr  *repository.ResponseRepository
	v   *Validator
}

func NewFormService(cfg config.Config, db *database.Database) *FormService {
	return &FormService{
		cfg: cfg,
		db:  db,
		fr:  repository.NewFormRepository(db),
		rr:  repository.NewResponseRepository(db),
		v:   NewValidator(),
	}
}

func (f *FormService) Create(formRequest *models.FormRequest, user *models.User) (*models.Form, types.ServiceErrors) {
	if err := f.v.Validate(formRequest); err != nil {
		return nil, err
	}

	form := formRequest.ToForm(user.ID)

	if err := f.fr.Create(&form); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, types.SvcError("Form code already exists", err)
		}

		return nil, types.SvcError("Failed to create form", err)
	}

	return &form, nil
}

func (f *FormService) GetForUser(user *models.User) ([]models.Form, error) {
	forms, err := f.fr.GetForUser(user.ID.String())
	if err != nil {
		return nil, types.SvcError("Error finding forms", err)
	}

	return forms, nil
}

func (f *FormService) GetByCode(code string, user *models.User) (*models.Form, error) {
	form, err := f.fr.GetByCode(code)
	if err != nil {
		return nil, types.SvcError("Form not found", err)
	}

	if !form.Multiple && form.UserID != user.ID {
		_, err := f.rr.GetByUserAndFormID(user.ID.String(), form.ID.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return form, nil
			}
		}

		return nil, types.SvcError("Form already submitted")
	}

	return form, nil
}

func (f *FormService) CheckCodeAvailability(code string) bool {
	_, err := f.fr.GetByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}
	return false
}

func (f *FormService) GetWithResponses(formID string, user *models.User) (*models.Form, types.ServiceErrors) {
	formWithResponses, err := f.fr.GetWithResponses(user.ID.String(), formID)
	if err != nil {
		return nil, types.SvcError("Error finding form", err)
	}

	return formWithResponses, nil
}
