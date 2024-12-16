package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/devmegablaster/bashform/internal/ui/form"
	"github.com/spf13/cobra"
)

func (c *CLI) FillForm() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "form [code]",
		Short:        "Fill out a form using form code",
		Args:         cobra.ExactArgs(1),
		Aliases:      []string{"f"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			client := services.NewClient(c.cfg.Api.BaseUrl, c.PubKey)
			formData, err := client.GetForm(args[0])
			if err != nil {
				cmd.Println(styles.Error.Render(constants.MessageFormNotFound))
				return nil
			}

			model := form.NewModel(formData, client, c.Session)
			p := tea.NewProgram(model,
				tea.WithInput(cmd.InOrStdin()),
				tea.WithOutput(cmd.OutOrStdout()),
			)
			_, err = p.Run()
			return err
		},
	}

	return formCmd
}
