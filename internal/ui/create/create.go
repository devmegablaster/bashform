package create

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/config"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
)

type Model struct {
	width, height int
	cfg           config.Config
	client        *services.Client
	code          string
	n             int
	questionsForm *huh.Form
	formResp      *models.FormResponse

	isCreating bool
	isCreated  bool
	err        error
	sizeError  bool
	init       bool
	initTime   time.Time
}

func NewModel(cfg config.Config, session ssh.Session, n int, code string, client *services.Client) *Model {
	pty, _, _ := session.Pty()

	sizeErr := false
	if pty.Window.Width < 50 || pty.Window.Height < 30 {
		sizeErr = true
	}

	return &Model{
		width:         pty.Window.Width,
		height:        pty.Window.Height,
		cfg:           cfg,
		code:          code,
		n:             n,
		client:        client,
		questionsForm: starterForm(n),

		sizeError: sizeErr,
		init:      true,
		initTime:  time.Now(),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.questionsForm.Init()
}

func (m *Model) View() string {
	if m.sizeError {
		return styles.Error.Render(fmt.Sprintf(constants.MessageSizeError, 50, 30, m.width, m.height))
	}

	var content string

	content = m.questionsForm.View()

	if m.isCreating {
		content = styles.Succes.Render(constants.MessageFormCreating)
	}

	if m.isCreated {
		content = styles.Succes.Render(constants.MessageFormCreated) +
			"\n\n" +
			styles.Description.Render(constants.MessageCommandHeader) +
			"\n" +
			styles.Heading.Render(fmt.Sprintf(constants.MessageCommand, m.cfg.SSH.URL, m.formResp.Data.Code)) +
			"\n\n" +
			styles.Description.Render(constants.MessageHelpExit)
	}

	if m.err != nil {
		content = styles.Error.Render(fmt.Sprintf(constants.MessageFormCreateError, m.err.Error()))
	}

	if m.init {
		return styles.PlaceCenter(m.width, m.height, constants.Logo)
	}

	box := styles.Box(m.width, content)

	return styles.PlaceCenterVertical(m.width,
		m.height,
		styles.Heading.MarginBottom(1).Render("New Form"),
		box,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Handle key presses
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.questionsForm.State == huh.StateCompleted {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	form, cmd := m.questionsForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.questionsForm = f
	}

	if m.questionsForm.State == huh.StateCompleted && !m.isCreating {
		m.CreateRequest()
	}

	if m.init && time.Since(m.initTime) > 2*time.Second {
		m.init = false
	}

	return m, cmd
}

func (m *Model) CreateRequest() {
	m.isCreating = true

	formRequest := huhToForm(m.n, m.questionsForm)
	formRequest.Code = m.code

	form, err := m.client.CreateForm(*formRequest)
	if err != nil {
		m.err = err
		return
	}

	m.formResp = form

	m.isCreated = true
}
