package models

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FormID    uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Answers   []Answer  `gorm:"foreignKey:ResponseID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type ResponseRequest struct {
	Answers []Answer `validate:"required"`
}

func (r *ResponseRequest) ToResponse(formID uuid.UUID, userID uuid.UUID) Response {
	return Response{
		UserID:  userID,
		FormID:  formID,
		Answers: r.Answers,
	}
}

func (r *Response) ToCSV() []string {
	csv := []string{}
	for _, answer := range r.Answers {
		csv = append(csv, answer.Value)
	}
	return csv
}
