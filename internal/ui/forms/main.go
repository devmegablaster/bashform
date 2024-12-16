package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
)

type sessionState int

const (
	formsView sessionState = iota
	responsesView
	responseView
)

type responseMsg struct {
	questions, answers []string
	formID             string
}

type responsesMsg struct {
	formID string
}

type formsMsg struct{}

func toResponses(formID string) tea.Cmd {
	return func() tea.Msg {
		return responsesMsg{formID}
	}
}

func toForms() tea.Cmd {
	return func() tea.Msg {
		return formsMsg{}
	}
}

func toResponse(questions, answers []string, formID string) tea.Cmd {
	return func() tea.Msg {
		return responseMsg{
			questions: questions,
			answers:   answers,
			formID:    formID,
		}
	}
}

type Model struct {
	width, height int
	state         sessionState
	formsModel    *formsModel
	rsm           *responsesModel
	rm            *ResponseModel
	spinner       spinner.Model
	sizeError     bool
}

func NewModel(items []models.Item, client *services.Client, session ssh.Session) *Model {
	p, _, _ := session.Pty()

	sizeErr := false
	if p.Window.Width < 50 || p.Window.Height < 30 {
		sizeErr = true
	}

	return &Model{
		width:      p.Window.Width,
		height:     p.Window.Height,
		state:      formsView,
		formsModel: newFormsModel(items, client, session),
		rsm:        newResponsesModel(client, session),
		rm:         NewResponseModel(session),
		sizeError:  sizeErr,
	}
}

func (m *Model) Init() tea.Cmd {

	return tea.Batch(m.formsModel.Init(), m.spinner.Tick, m.rsm.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case responsesMsg:
		m.rsm.formID = msg.formID
		m.rsm.GetResponses()
		m.state = responsesView
	case formsMsg:
		m.state = formsView
	case responseMsg:
		m.rm.questions = msg.questions
		m.rm.answers = msg.answers
		m.rm.formID = msg.formID
		m.state = responseView
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case formsView:
		newModel, cmd := m.formsModel.Update(msg)
		formsModel, ok := newModel.(*formsModel)
		if !ok {
			return m, cmd
		}
		m.formsModel = formsModel
		cmds = append(cmds, cmd)

	case responsesView:
		newModel, cmd := m.rsm.Update(msg)
		rsm, ok := newModel.(*responsesModel)
		if !ok {
			return m, cmd
		}
		m.rsm = rsm
		cmds = append(cmds, cmd)

	case responseView:
		newModel, cmd := m.rm.Update(msg)
		rm, ok := newModel.(*ResponseModel)
		if !ok {
			return m, cmd
		}
		m.rm = rm
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.sizeError {
		return styles.Error.Render(fmt.Sprintf(constants.MessageSizeError, 50, 30, m.width, m.height))
	}

	switch m.state {
	case formsView:
		return m.formsModel.View()
	case responsesView:
		return m.rsm.View()
	case responseView:
		return m.rm.View()
	}
	return ""
}
