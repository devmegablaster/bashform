package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/devmegablaster/bashform/internal/ui/forms"
	"github.com/spf13/cobra"
)

func (c *CLI) getForms() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "forms",
		Short:        "Get your forms",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			userForms, err := c.formSvc.GetForUser(c.user)
			if err != nil {
				cmd.Println(styles.Error.Render(err.Error()))
				return nil
			}

			items := models.FormsToItems(userForms)
			model := forms.NewModel(items, c.formSvc, c.Session)

			p := tea.NewProgram(model,
				tea.WithAltScreen(),
				tea.WithInput(c.Session),
				tea.WithOutput(c.Session),
			)

			_, err = p.Run()

			return err
		},
	}

	return formCmd
}
