package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/devmegablaster/bashform/internal/ui/form"
	"github.com/spf13/cobra"
)

func (c *CLI) fillForm() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "form [code]",
		Short:        "Fill out a form using form code",
		Args:         cobra.ExactArgs(1),
		Aliases:      []string{"f"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			formData, err := c.formSvc.GetByCode(args[0], c.user)
			if err != nil {
				c.logger.Error("Error getting form", "error", err)
				cmd.Println(styles.Error.Render(err.Error()))
				return nil
			}

			c.logger.Info("Form retrieved", "form", formData.ID)

			model := form.NewModel(formData, c.responseSvc, c.Session)
			p := tea.NewProgram(model,
				tea.WithAltScreen(),
				tea.WithInput(cmd.InOrStdin()),
				tea.WithOutput(cmd.OutOrStdout()),
			)
			_, err = p.Run()
			return err
		},
	}

	return formCmd
}
