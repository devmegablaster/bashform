package models

import (
	"time"

	"github.com/charmbracelet/huh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/google/uuid"
)

type Form struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:varchar(255);"`
	Code        string     `gorm:"type:varchar(20);not null;unique"`
	Questions   []Question `gorm:"foreignKey:FormID"`
	Responses   []Response `gorm:"foreignKey:FormID"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	Multiple    bool       `gorm:"type:boolean;not null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

type FormRequest struct {
	Name        string `validate:"required"`
	Description string
	Questions   []QuestionRequest `validate:"required"`
	Code        string            `validate:"required"`
	Multiple    bool
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

func (f *Form) ToItem() Item {
	return Item{
		ID:   f.ID.String(),
		Name: f.Name,
		Desc: f.Code,
	}
}

func FormsToItems(forms []Form) []Item {
	items := []Item{}
	for _, form := range forms {
		items = append(items, form.ToItem())
	}
	return items
}
