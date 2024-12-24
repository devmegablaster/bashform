package models

import (
	"github.com/charmbracelet/huh"
	"github.com/devmegablaster/bashform/internal/constants"
)

type Form struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Code        string     `json:"code"`
	Questions   []Question `json:"questions"`
	Responses   []Response `json:"responses"`
	Multiple    bool       `json:"multiple"`
	Error       string     `json:"error"`
}

type FormResponse struct {
	Data  Form   `json:"data"`
	Error string `json:"error"`
}

type FormRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Code        string            `json:"code"`
	Questions   []QuestionRequest `json:"questions"`
	Multiple    bool              `json:"multiple"`
}

func (f Form) ToHuhForm() *huh.Form {
	fields := []huh.Field{}

	for _, question := range f.Questions {
		switch question.Type {
		case constants.FIELD_TEXT:
			field := huh.NewInput().Title(question.Text).Key(question.ID)
			if question.Required {
				field = field.Validate(huh.ValidateNotEmpty())
			}
			fields = append(fields, field)

		case constants.FIELD_TEXTAREA:
			field := huh.NewText().Title(question.Text).Key(question.ID)
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

			field := huh.NewSelect[string]().Title(question.Text).Options(opts...).Key(question.ID)
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
