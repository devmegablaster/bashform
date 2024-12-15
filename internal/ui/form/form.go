package form

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
)

type Model struct {
	width, height int
	Form          models.Form
	session       ssh.Session
	client        *services.Client
	huhForm       *huh.Form
	spinner       spinner.Model

	isSubmitting  bool
	sizeError     bool
	submitError   error
	submitSuccess bool
	init          bool
	initTime      time.Time
}

func NewModel(form models.Form, client *services.Client, session ssh.Session) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Succes

	pty, _, _ := session.Pty()

	sizeErr := false
	if pty.Window.Width < 50 || pty.Window.Height < 35 {
		sizeErr = true
	}

	return &Model{
		Form:      form,
		huhForm:   form.ToHuhForm(),
		spinner:   s,
		client:    client,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		session:   session,
		init:      true,
		sizeError: sizeErr,
		initTime:  time.Now(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.huhForm.Init(), tea.EnterAltScreen)
}

func (m *Model) View() string {
	var content string

	content = m.huhForm.View()

	if m.isSubmitting {
		content = m.spinner.View() + "\n" + styles.Description.Render(constants.MessageFormSubmitting)
	}

	if m.submitSuccess {
		content = styles.Succes.Render(constants.MessageFormSubmitted) + "\n\n" + styles.Description.Render(constants.MessageHelpExit)
	}

	if m.init {
		return styles.PlaceCenter(m.width, m.height, constants.Logo)
	}

	if m.sizeError {
		return styles.Error.Render(fmt.Sprintf(constants.MessageSizeError, m.width, m.height))
	}

	box := styles.Box(m.width, content)

	return styles.PlaceCenterVertical(m.width,
		m.height,
		styles.Heading.Render(m.Form.Name),
		styles.Description.MarginBottom(1).Render(m.Form.Description),
		box,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Handle key presses for exit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "q" {
			if m.huhForm.State == huh.StateCompleted {
				return m, tea.Quit
			}
		}
	}

	cmds := []tea.Cmd{}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	form, cmd := m.huhForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.huhForm = f
	}

	if m.huhForm.State == huh.StateCompleted && !m.isSubmitting {
		m.Submit()
	}

	if time.Since(m.initTime) > 2*time.Second {
		m.init = false
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Submit() {
	m.isSubmitting = true
	answer := []models.Answer{}

	for _, question := range m.Form.Questions {
		answer = append(answer, models.Answer{
			QuestionID: question.ID,
			Value:      m.huhForm.GetString(question.ID),
		})
	}

	response := models.Response{
		FormID:  m.Form.ID,
		Answers: answer,
	}

	m.client.SubmitForm(m.Form.Code, response)

	m.submitSuccess = true
}