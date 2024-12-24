package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/services"
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
			client := services.NewClient(c.cfg.Api.BaseURL, c.PubKey)
			formData, err := client.GetForm(args[0])
			if err != nil {
				cmd.Println(styles.Error.Render(err.Error()))
				return nil
			}

			model := form.NewModel(formData, client, c.Session)
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
