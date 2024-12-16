package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/devmegablaster/bashform/internal/ui/forms"
	"github.com/spf13/cobra"
)

func (c *CLI) GetForms() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "forms",
		Short:        "Get your forms",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			client := services.NewClient(c.cfg.Api.BaseUrl, c.PubKey)
			formsResp, err := client.GetForms()
			if err != nil {
				cmd.Println(styles.Error.Render(constants.MessageFormNotFound))
				return nil
			}

			items := models.FormsToItems(formsResp)
			model := forms.NewModel(items, client, c.Session)

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
