package models

import (
	"time"

	"github.com/charmbracelet/huh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/google/uuid"
)

type Form struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:varchar(255);" json:"description"`
	Code        string     `gorm:"type:varchar(20);not null;unique" json:"code"`
	Questions   []Question `gorm:"foreignKey:FormID" json:"questions"`
	Responses   []Response `gorm:"foreignKey:FormID" json:"responses"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Multiple    bool       `gorm:"type:boolean;not null" json:"multiple"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

type FormRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Questions   []QuestionRequest `json:"questions" validate:"required"`
	Code        string            `json:"code" validate:"min=1,max=20"`
	Multiple    bool              `json:"multiple"`
}

type FormResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Code        string             `json:"code"`
	Questions   []QuestionResponse `json:"questions"`
	UserID      uuid.UUID          `json:"user_id"`
	Multiple    bool               `json:"multiple"`
}

func (f *Form) ToResponse() FormResponse {
	questions := make([]QuestionResponse, len(f.Questions))
	for i, q := range f.Questions {
		questions[i] = q.ToResponse()
	}

	return FormResponse{
		ID:        f.ID,
		Name:      f.Name,
		Code:      f.Code,
		Questions: questions,
		UserID:    f.UserID,
		Multiple:  f.Multiple,
	}
}

func (f *FormRequest) ToForm(userID uuid.UUID) Form {
	questions := make([]Question, len(f.Questions))
	for i, q := range f.Questions {
		questions[i] = q.ToQuestion()
	}

	return Form{
		Name:        f.Name,
		Description: f.Description,
		Questions:   questions,
		Code:        f.Code,
		UserID:      userID,
		Multiple:    f.Multiple,
	}
}

func (f Form) ToHuhForm() *huh.Form {
	fields := []huh.Field{}

	for _, question := range f.Questions {
		switch question.Type {
		case constants.FIELD_TEXT:
			field := huh.NewInput().Title(question.Text).Key(question.ID.String())
			if question.Required {
				field = field.Validate(huh.ValidateNotEmpty())
			}
			fields = append(fields, field)

		case constants.FIELD_TEXTAREA:
			field := huh.NewText().Title(question.Text).Key(question.ID.String())
			if question.Required {
				field = field.Validate(huh.ValidateNotEmpty())
			}
			fields = append(fields, field)

		case constants.FIELD_SELECT:
			options := make([]string, len(question.Options))
			for i, option := range question.Options {
				options[i] = option.Text
			}

			opts := huh.NewOptions(options...)

			field := huh.NewSelect[string]().Title(question.Text).Options(opts...).Key(question.ID.String())
			if question.Required {
				field = field.Validate(huh.ValidateNotEmpty())
			}
			fields = append(fields, field)
		}
	}

	rootGroup := huh.NewGroup(fields...).WithTheme(huh.ThemeCharm())

	form := huh.NewForm(rootGroup).WithTheme(huh.ThemeCharm())
	return form
}
