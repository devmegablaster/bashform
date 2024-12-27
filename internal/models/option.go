package models

import "github.com/google/uuid"

type Option struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	Text       string    `gorm:"type:varchar(255);not null"`
}

type OptionRequest struct {
	Text string
}

func (o *OptionRequest) ToOption() Option {
	return Option{
		Text: o.Text,
	}
}
