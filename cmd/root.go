package cmd

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/database"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/spf13/cobra"
)

type CLI struct {
	cfg      config.Config
	database *database.Database
	Session  ssh.Session
	user     *models.User
	RootCmd  *cobra.Command
	logger   *slog.Logger

	formSvc     *services.FormService
	responseSvc *services.ResponseService
}

func NewCLI(cfg config.Config, db *database.Database, session ssh.Session) *CLI {
	rootCmd := &cobra.Command{
		Use: "bashform",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(constants.Logo)
			cmd.Help()
			return nil
		},
	}

	u := session.Context().Value("user").(*models.User)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})).With("user", u.ID).With("username", u.Name)

	c := &CLI{
		cfg:         cfg,
		Session:     session,
		user:        u,
		RootCmd:     rootCmd,
		logger:      logger,
		formSvc:     services.NewFormService(cfg, db, logger),
		responseSvc: services.NewResponseService(cfg, db, logger),
	}

	c.Init()
	return c
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
	c.AddCommand(c.fillForm())
	c.AddCommand(c.createForm())
	c.AddCommand(c.getForms())
}

func (c *CLI) Run() error {
	return c.RootCmd.Execute()
}
