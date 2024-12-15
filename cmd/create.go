package cmd

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devmegablaster/bashform/internal/create"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/spf13/cobra"
)

func (c *CLI) Create() *cobra.Command {
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
				return fmt.Errorf("number of questions must be greater than 0")
			}

			var code string

			if len(args) < 2 {
				code = ""
			} else {
				code = args[1]
			}

			client := services.NewClient(c.cfg.Api.BaseUrl, c.PubKey)

			check, err := client.GetForm(code)

			if err != nil {
				cmd.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render("\nError checking code...\n"))
			}

			if check.Code != "" {
				cmd.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render("\nA Form with that code already exists!\n"))
				return nil
			}

			cr := create.NewModel(c.Session, n, code, client)

			p := tea.NewProgram(cr,
				tea.WithAltScreen(),
				tea.WithInput(c.Session),
				tea.WithOutput(c.Session))

			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
				os.Exit(1)
			}

			return nil
		},
	}

	return newFormCmd
}
