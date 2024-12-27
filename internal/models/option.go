package models

import "github.com/google/uuid"

type Option struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null" json:"question_id"`
	Text       string    `gorm:"type:varchar(255);not null" json:"text"`
}

type OptionRequest struct {
	Text string `json:"text"`
}

type OptionResponse struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Text       string    `json:"text"`
}

func (o *Option) ToResponse() OptionResponse {
	return OptionResponse{
		ID:         o.ID,
		QuestionID: o.QuestionID,
		Text:       o.Text,
	}
}

func (o *OptionRequest) ToOption() Option {
	return Option{
		Text: o.Text,
	}
}
