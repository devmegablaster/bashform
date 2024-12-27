package services

import (
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/repository"
	"github.com/devmegablaster/bashform/internal/types"
	"github.com/google/uuid"
)

type ResponseService struct {
	cfg config.Config
	db  *database.Database
	rr  *repository.ResponseRepository
	fr  *repository.FormRepository
	v   Validator
}

func NewResponseService(cfg config.Config, db *database.Database) *ResponseService {
	return &ResponseService{
		cfg: cfg,
		db:  db,
		rr:  repository.NewResponseRepository(db),
		fr:  repository.NewFormRepository(db),
		v:   *NewValidator(),
	}
}

func (r *ResponseService) Create(responseRequest models.ResponseRequest, formID uuid.UUID, user *models.User) (models.Response, types.ServiceErrors) {
	if err := r.v.Validate(responseRequest); err != nil {
		return models.Response{}, err
	}
	response := responseRequest.ToResponse(formID, user.ID)

	_, err := r.rr.GetByUserAndFormID(user.ID.String(), response.FormID.String())
	if err == nil {
		return models.Response{}, types.SvcError("Response already exists")
	}

	if err := r.rr.Create(&response); err != nil {
		return models.Response{}, types.SvcError("Failed to create response", err)
	}

	return response, nil
}

func (r *ResponseService) GetByFormID(formID string, user *models.User) (*models.Form, types.ServiceErrors) {
	responses, err := r.fr.GetWithResponses(user.ID.String(), formID)
	if err != nil {
		return nil, types.SvcError("Error finding form", err)
	}

	return responses, nil
}
