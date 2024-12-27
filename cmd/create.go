package cmd

import (
	"log/slog"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/styles"
	"github.com/devmegablaster/bashform/internal/ui/create"
	"github.com/spf13/cobra"
)

func (c *CLI) createForm() *cobra.Command {
	newFormCmd := &cobra.Command{
		Use:     "create [number of questions] [share code]",
		Args:    cobra.ExactArgs(2),
		Short:   "Create a new form with a specific number of questions and shareable code",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			if n < 1 {
				cmd.Println(styles.Error.Render("number of questions must be greater than 0"))
				return nil
			}

			var code string

			if len(args) < 2 {
				code = ""
			} else {
				code = args[1]
			}

			avl := c.formSvc.CheckCodeAvailability(code)

			if !avl {
				cmd.Println(styles.Error.Render(constants.MessageCodeNotAvailable))
				return nil
			}

			cr := create.NewModel(c.cfg, c.Session, n, code, c.formSvc)

			p := tea.NewProgram(cr,
				tea.WithAltScreen(),
				tea.WithInput(c.Session),
				tea.WithOutput(c.Session))

			if _, err := p.Run(); err != nil {
				slog.Error("Error running program", "error", err)
				os.Exit(1)
			}

			return nil
		},
	}

	return newFormCmd
}
