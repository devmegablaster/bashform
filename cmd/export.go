package cmd

import (
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/spf13/cobra"
)

func (c *CLI) exportResponses() *cobra.Command {
	formCmd := &cobra.Command{
		Use:          "export [code]",
		Short:        "Exports responses for a form using form code to CSV",
		Args:         cobra.ExactArgs(1),
		Aliases:      []string{"e"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			formCSV, err := c.formSvc.GetResponsesCSV(args[0], c.user)
			if err != nil {
				c.logger.Error("Error getting form", "error", err)
				cmd.Println(styles.Error.Render(err.Error()))
			}

			c.logger.Info("Form exported to CSV")

			cmd.Println(formCSV)
			return nil
		},
	}

	return formCmd
}
