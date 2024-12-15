package form

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
)

type Model struct {
	width         int
	height        int
	Form          models.Form
	client        *services.Client
	huhForm       *huh.Form
	spinner       spinner.Model
	isSubmitting  bool
	submitError   error
	submitSuccess bool
	init          bool
	initTime      time.Time
}

func NewModel(form models.Form, client *services.Client, session ssh.Session) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))

	pty, _, _ := session.Pty()

	return &Model{
		Form:     form,
		huhForm:  form.ToHuhForm(),
		spinner:  s,
		client:   client,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		init:     true,
		initTime: time.Now(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.huhForm.Init(), tea.EnterAltScreen)
}

func (m *Model) View() string {
	var content string

	content = m.huhForm.View()

	if m.isSubmitting {
		content = m.spinner.View() + "\n Submitting your response..."
	}

	if m.submitSuccess {
		successMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("Form submitted successfully!")
		helpMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("q or ctrl+c to exit")

		content = successMsg + "\n\n" + helpMsg
	}

	if m.init {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, constants.Logo)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#94a3b8")).
		Align(lipgloss.Center).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Underline(true).Render(m.Form.Name),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b")).MarginBottom(1).Render(m.Form.Description),
			box,
		))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
