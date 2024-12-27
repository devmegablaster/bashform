package services

import (
	"fmt"
	"log/slog"

	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
	"github.com/google/uuid"
)

type ResponseService struct {
	cfg    config.Config
	db     *database.Database
	rr     *repository.ResponseRepository
	fr     *repository.FormRepository
	v      Validator
	logger *slog.Logger
}

func NewResponseService(cfg config.Config, db *database.Database, logger *slog.Logger) *ResponseService {
	return &ResponseService{
		cfg:    cfg,
		db:     db,
		rr:     repository.NewResponseRepository(db),
		fr:     repository.NewFormRepository(db),
		v:      *NewValidator(),
		logger: logger,
	}
}

func (r *ResponseService) Create(responseRequest models.ResponseRequest, formID uuid.UUID, user *models.User) (models.Response, error) {
	// TODO: implement encryption of responses
	if err := r.v.Validate(responseRequest); err != nil {
		return models.Response{}, err
	}
	response := responseRequest.ToResponse(formID, user.ID)

	_, err := r.rr.GetByUserAndFormID(user.ID.String(), response.FormID.String())
	if err == nil {
		r.logger.Warn("Response already exists for form", slog.String("form_id", response.FormID.String()))
		return models.Response{}, fmt.Errorf("Response already exists for form")
	}

	if err := r.rr.Create(&response); err != nil {
		r.logger.Error("Failed to create response", "error", err)
		return models.Response{}, fmt.Errorf("Failed to create response")
	}

	return response, nil
}

func (r *ResponseService) GetByFormID(formID string, user *models.User) (*models.Form, error) {
	responses, err := r.fr.GetWithResponses(user.ID.String(), formID)
	if err != nil {
		r.logger.Error("Error finding form", slog.String("form_id", formID), "error", err)
		return nil, fmt.Errorf("Error finding form")
	}

	return responses, nil
}
