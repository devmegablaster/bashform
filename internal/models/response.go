package models

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FormID    uuid.UUID `gorm:"type:uuid;not null" json:"form_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Answers   []Answer  `gorm:"foreignKey:ResponseID" json:"answers"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type ResponseResponse struct {
	ID        uuid.UUID  `json:"id"`
	FormID    uuid.UUID  `json:"form_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Answers   []Answer   `json:"answers"`
	Questions []Question `json:"questions"`
}

type ResponseRequest struct {
	Answers []Answer `json:"answers" validate:"required"`
}

func (r *ResponseRequest) ToResponse(formID uuid.UUID, userID uuid.UUID) Response {
	return Response{
		UserID:  userID,
		FormID:  formID,
		Answers: r.Answers,
	}
}
