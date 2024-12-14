package models

import (
	"github.com/charmbracelet/huh"
)

type FormResponse struct {
	Data Form `json:"data"`
}

type Form struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Code        string     `json:"code"`
	Questions   []Question `json:"questions"`
}

type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (f Form) ToHuhForm() *huh.Form {
	fields := []huh.Field{}

	for _, question := range f.Questions {
		switch question.Type {
		case "text":
			fields = append(fields, huh.NewInput().Title(question.Text).Key(question.ID))
		case "textarea":
			fields = append(fields, huh.NewText().Title(question.Text).Key(question.ID))
		case "select":
			options := make([]string, len(question.Options))
			for i, option := range question.Options {
				options[i] = option.Text
			}

			opts := huh.NewOptions(options...)
			fields = append(fields, huh.NewSelect[string]().Title(question.Text).Options(opts...).Key(question.ID))
		}
	}

	rootGroup := huh.NewGroup(fields...).WithTheme(huh.ThemeCharm())

	form := huh.NewForm(rootGroup).WithTheme(huh.ThemeCharm())
	return form
}
