package services

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
	"gorm.io/gorm"
)

type FormService struct {
	cfg    config.Config
	db     *database.Database
	fr     *repository.FormRepository
	rr     *repository.ResponseRepository
	v      *Validator
	logger *slog.Logger
}

func NewFormService(cfg config.Config, db *database.Database, logger *slog.Logger) *FormService {
	return &FormService{
		cfg:    cfg,
		db:     db,
		fr:     repository.NewFormRepository(db),
		rr:     repository.NewResponseRepository(db),
		v:      NewValidator(),
		logger: logger,
	}
}

// create a new form using the form request for the user
func (f *FormService) Create(formRequest *models.FormRequest, user *models.User) (*models.Form, error) {
	if err := f.v.Validate(formRequest); err != nil {
		f.logger.Error("Failed to validate form request", "error", err)
		return nil, err
	}

	form := formRequest.ToForm(user.ID)

	if err := f.fr.Create(&form); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			f.logger.Warn("Form with code already exists", slog.String("code", form.Code))
			return nil, fmt.Errorf("Form with code %s already exists", form.Code)
		}

		f.logger.Error("Failed to create form", "error", err)
		return nil, fmt.Errorf("Failed to create form")
	}

	return &form, nil
}

// get all forms for the user
func (f *FormService) GetForUser(user *models.User) ([]models.Form, error) {
	forms, err := f.fr.GetForUser(user.ID.String())
	if err != nil {
		f.logger.Error("Error finding forms", "error", err)
		return nil, fmt.Errorf("Error finding forms")
	}

	return forms, nil
}

// get a form by its code
func (f *FormService) GetByCode(code string, user *models.User) (*models.Form, error) {
	form, err := f.fr.GetByCode(code)
	if err != nil {
		f.logger.Error("Error finding form", slog.String("code", code), "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Form not found")
		}
		return nil, fmt.Errorf("Error finding form")
	}

	if !form.Multiple && form.UserID != user.ID {
		_, err := f.rr.GetByUserAndFormID(user.ID.String(), form.ID.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return form, nil
			}
		}

		f.logger.Warn("Form already submitted", slog.String("form_id", form.ID.String()))
		return nil, fmt.Errorf("Form already submitted")
	}

	return form, nil
}

// check if the form code is available
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

// get user owned form with its responses
func (f *FormService) GetWithResponses(formID string, user *models.User) (*models.Form, error) {
	// TODO: Implement decryption of responses
	formWithResponses, err := f.fr.GetWithResponses(user.ID.String(), formID)
	if err != nil {
		f.logger.Error("Error finding form", slog.String("form_id", formID), "error", err)
		return nil, fmt.Errorf("Error finding form")
	}

	return formWithResponses, nil
}
