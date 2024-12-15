package create

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/devmegablaster/bashform/internal/models"
)

func questionForm(index int) *huh.Group {
	return huh.NewGroup(
		huh.NewSelect[string]().
			Title("Question Type").
			Options(
				huh.NewOption("Text", "text"),
				huh.NewOption("Textarea", "textarea"),
				huh.NewOption("Select", "select"),
			).Key(fmt.Sprintf("type%d", index)).Validate(huh.ValidateNotEmpty()),
		huh.NewInput().Title("Question").Key(fmt.Sprintf("question%d", index)).Validate(huh.ValidateNotEmpty()),
		huh.NewInput().Title("Options (comma separated) [For select questions only]").Key(fmt.Sprintf("options%d", index)),
		huh.NewConfirm().Title("Required?").Key(fmt.Sprintf("required%d", index)),
	)
}

func starterForm(n int) *huh.Form {
	questions := make([]*huh.Group, n)
	for i := 0; i < n; i++ {
		questions[i] = questionForm(i)
	}

	questions = append(questions, huh.NewGroup(
		huh.NewConfirm().Title("Allow multiple submissions by a user?").Key("allowMultiple"),
	))

	allGroups := append([]*huh.Group{
		huh.NewGroup(
			huh.NewInput().Title("Form Name").Key("name").Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Form Description").Key("description").Validate(huh.ValidateNotEmpty()),
		),
	}, questions...)

	return huh.NewForm(
		allGroups...,
	)
}

func huhToForm(n int, huhForm *huh.Form) *models.FormRequest {
	questions := []models.QuestionRequest{}

	for i := 0; i < n; i++ {
		question := models.QuestionRequest{
			Text: huhForm.GetString(fmt.Sprintf("question%d", i)),
			Type: huhForm.GetString(fmt.Sprintf("type%d", i)),
			// Required: huhForm.GetBool(fmt.Sprintf("required%d", i)),
		}

		if question.Type == "select" {
			optsStr := huhForm.GetString(fmt.Sprintf("options%d", i))
			optionRequests := []models.OptionRequest{}

			opts := strings.Split(optsStr, ",")
			for _, opt := range opts {
				optionRequests = append(optionRequests, models.OptionRequest{Text: strings.TrimSpace(opt)})
			}

			question.Options = optionRequests
		}

		questions = append(questions, question)
	}

	formRequest := models.FormRequest{
		Name:        huhForm.GetString("name"),
		Description: huhForm.GetString("description"),
		Questions:   questions,
		// AllowMultiple: huhForm.GetBool("allowMultiple"),
	}

	return &formRequest
}
