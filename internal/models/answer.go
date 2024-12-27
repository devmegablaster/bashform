package models

import "github.com/google/uuid"

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null" json:"question_id"`
	Value      string    `gorm:"type:text" json:"value"`
	ResponseID uuid.UUID `gorm:"type:uuid;not null" json:"response_id"`
}

type AnswerRequest struct {
	QuestionID uuid.UUID `json:"question_id" validate:"required"`
	Value      string    `json:"value" validate:"required"`
}
