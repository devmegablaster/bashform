package forms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/styles"
)

type ResponseModel struct {
	width, height      int
	formID             string
	session            ssh.Session
	questions, answers []string
}

func NewResponseModel(session ssh.Session) *ResponseModel {
	p, _, _ := session.Pty()

	return &ResponseModel{
		width:     p.Window.Width,
		height:    p.Window.Height,
		session:   session,
		questions: []string{},
		answers:   []string{},
		formID:    "",
	}
}

func (m *ResponseModel) Init() tea.Cmd {
	return nil
}

func (m *ResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return nil, toResponses(m.formID)
		}
	}

	return m, nil
}

func (m *ResponseModel) View() string {
	var content string
	for i, question := range m.questions {
		content += styles.Heading.Render(question)
		content += "\n"
		content += styles.Description.Render(m.answers[i])
		content += "\n\n"
	}

	content += styles.Description.Render("\n Press ESC to go back")

	return styles.PlaceCenterVertical(m.width, m.height, styles.Heading.Render("Response\n"), styles.Box(m.width, content))
}
