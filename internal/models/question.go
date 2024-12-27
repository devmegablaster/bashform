package models

import "github.com/google/uuid"

type Question struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FormID   uuid.UUID `gorm:"type:uuid;not null"`
	Text     string    `gorm:"type:varchar(255);not null"`
	Type     string    `gorm:"type:varchar(255);not null"`
	Options  []Option  `gorm:"foreignKey:QuestionID"`
	Required bool      `gorm:"type:boolean;not null"`
}

type QuestionRequest struct {
	Text     string `validate:"required,max=255,min=1"`
	Type     string `validate:"required,max=255,min=1"`
	Options  []OptionRequest
	Required bool `validate:"required"`
}

func (q *QuestionRequest) ToQuestion() Question {
	options := make([]Option, len(q.Options))
	for i, o := range q.Options {
		options[i] = o.ToOption()
	}

	return Question{
		Text:     q.Text,
		Type:     q.Type,
		Options:  options,
		Required: q.Required,
	}
}
