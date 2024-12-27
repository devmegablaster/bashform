package models

import "github.com/google/uuid"

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	Value      string    `gorm:"type:text"`
	ResponseID uuid.UUID `gorm:"type:uuid;not null"`
}

type AnswerRequest struct {
	QuestionID uuid.UUID `validate:"required"`
	Value      string    `validate:"required"`
}
