package models

import "github.com/google/uuid"

type Question struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FormID   uuid.UUID `gorm:"type:uuid;not null" json:"form_id"`
	Text     string    `gorm:"type:varchar(255);not null" json:"text"`
	Type     string    `gorm:"type:varchar(255);not null" json:"type"`
	Options  []Option  `gorm:"foreignKey:QuestionID" json:"options"`
	Required bool      `gorm:"type:boolean;not null" json:"required"`
}

type QuestionRequest struct {
	Text     string          `json:"text" validate:"required,max=255,min=1"`
	Type     string          `json:"type" validate:"required,max=255,min=1"`
	Options  []OptionRequest `json:"options"`
	Required bool            `json:"required" validate:"required"`
}

type QuestionResponse struct {
	ID       uuid.UUID        `json:"id"`
	FormID   uuid.UUID        `json:"form_id"`
	Text     string           `json:"text"`
	Type     string           `json:"type"`
	Options  []OptionResponse `json:"options"`
	Required bool             `json:"required"`
}

func (q *Question) ToResponse() QuestionResponse {
	options := make([]OptionResponse, len(q.Options))
	for i, o := range q.Options {
		options[i] = o.ToResponse()
	}

	return QuestionResponse{
		ID:       q.ID,
		FormID:   q.FormID,
		Text:     q.Text,
		Type:     q.Type,
		Options:  options,
		Required: q.Required,
	}
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
