package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/form"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/spf13/cobra"
)

func (c *CLI) Form() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "form [code]",
		Short:        "Fill out a form",
		Args:         cobra.ExactArgs(1),
		Aliases:      []string{"f"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			client := services.NewClient(c.cfg.Api.BaseUrl, c.PubKey)
			formData := client.GetForm(args[0])

			model := form.NewModel(formData, client, c.Session)
			p := tea.NewProgram(model,
				tea.WithMouseCellMotion(),
				tea.WithInput(cmd.InOrStdin()),
				tea.WithOutput(cmd.OutOrStdout()),
			)
			_, err := p.Run()
			return err
		},
	}

	return formCmd
}
