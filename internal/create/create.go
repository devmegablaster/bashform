package create

import (
	"fmt"
	"strings"
	"time"

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
	questionsForm *huh.Form
	code          string
	n             int
	client        *services.Client
	isCreating    bool
	isCreated     bool
	err           error
	init          bool
	initTime      time.Time
	formResp      *models.FormResponse
}

func NewModel(session ssh.Session, n int, code string, client *services.Client) *Model {
	pty, _, _ := session.Pty()

	return &Model{
		width:         pty.Window.Width,
		height:        pty.Window.Height,
		questionsForm: starterForm(n),
		code:          code,
		n:             n,
		client:        client,
		init:          true,
		initTime:      time.Now(),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.questionsForm.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	return m, cmd
}

func (m *Model) View() string {

	var content string

	content = m.questionsForm.View()

	if m.isCreating {
		content = "Creating your form..."
	}

	if m.isCreated {
		successMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("Form created successfully!")
		linkMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6")).Render("Your Access Command:")
		accessMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6")).Render(fmt.Sprintf("ssh -t bashform.me f %s", m.formResp.Data.Code))
		helpMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("q or ctrl+c to exit")

		content = successMsg + "\n\n" + linkMsg + "\n" + accessMsg + "\n\n" + helpMsg
	}

	if m.err != nil {
		content = fmt.Sprintf("Error creating form: %s", m.err)
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
			lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6")).MarginBottom(1).Bold(true).Render("New Form"),
			box,
		))
}

func (m *Model) CreateRequest() {
	m.isCreating = true

	questions := []models.QuestionRequest{}

	for i := 0; i < m.n; i++ {
		question := models.QuestionRequest{
			Text: m.questionsForm.GetString(fmt.Sprintf("question%d", i)),
			Type: m.questionsForm.GetString(fmt.Sprintf("type%d", i)),
		}

		options := strings.Split(m.questionsForm.GetString(fmt.Sprintf("options%d", i)), ",")
		optionRequests := []models.OptionRequest{}
		for _, option := range options {
			optionRequests = append(optionRequests, models.OptionRequest{Text: strings.TrimSpace(option)})
		}

		question.Options = optionRequests
		questions = append(questions, question)
	}

	formRequest := models.FormRequest{
		Name:        m.questionsForm.GetString("name"),
		Description: m.questionsForm.GetString("description"),
		Code:        m.code,
		Questions:   questions,
	}

	form, err := m.client.CreateForm(formRequest)
	if err != nil {
		m.err = err
		return
	}

	m.formResp = form

	m.isCreated = true
}
