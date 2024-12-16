package cmd

import (
	"bytes"
	"encoding/base64"

	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/spf13/cobra"
)

type CLI struct {
	cfg     config.Config
	Session ssh.Session
	PubKey  string
	RootCmd *cobra.Command
}

func NewCLI(cfg config.Config, session ssh.Session) *CLI {

	encoded := bytes.Buffer{}

	enc := base64.NewEncoder(base64.StdEncoding, &encoded)
	enc.Write(session.PublicKey().Marshal())

	rootCmd := &cobra.Command{
		Use: "bashform",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	return &CLI{
		cfg:     cfg,
		Session: session,
		PubKey:  encoded.String(),
		RootCmd: rootCmd,
	}
}

func (c *CLI) AddCommand(cmd *cobra.Command) {
	c.RootCmd.AddCommand(cmd)
}

func (c *CLI) Init() {
	c.RootCmd.SetArgs(c.Session.Command())
	c.RootCmd.SetIn(c.Session)
	c.RootCmd.SetOut(c.Session)
	c.RootCmd.SetErr(c.Session.Stderr())
	c.RootCmd.SetContext(c.Session.Context())

	// Add Commands
	c.AddCommand(c.Form())
	c.AddCommand(c.Create())
	c.AddCommand(c.GetForms())
}

func (c *CLI) Run() error {
	return c.RootCmd.Execute()
}
