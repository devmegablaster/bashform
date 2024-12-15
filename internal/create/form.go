package create

import (
	"fmt"

	"github.com/charmbracelet/huh"
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
