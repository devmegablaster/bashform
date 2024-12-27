package forms

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/services"
	"github.com/devmegablaster/bashform/internal/styles"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type responsesModel struct {
	width, height int
	formID        string
	Form          models.Form
	session       ssh.Session
	user          *models.User
	formSvc       *services.FormService
	table         table.Model

	sizeError     bool
	responseError error
	init          bool
	initTime      time.Time
}

func newResponsesModel(formSvc *services.FormService, session ssh.Session) *responsesModel {
	pty, _, _ := session.Pty()

	sizeErr := false
	if pty.Window.Width < 50 || pty.Window.Height < 35 {
		sizeErr = true
	}

	t := table.New(
		table.WithFocused(true),
		table.WithHeight(20),
	)

	t.SetStyles(styles.TableStyle())

	u := session.Context().Value("user").(*models.User)

	return &responsesModel{
		Form:      models.Form{},
		formSvc:   formSvc,
		table:     t,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		session:   session,
		user:      u,
		init:      true,
		sizeError: sizeErr,
		initTime:  time.Now(),
	}
}

func (m *responsesModel) Init() tea.Cmd {
	return nil
}

func (m *responsesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return nil, toForms()
		}
		if msg.Type == tea.KeyEnter {
			questions := []string{}
			q := m.table.Columns()
			for _, column := range q {
				questions = append(questions, column.Title)
			}
			answers := m.table.SelectedRow()
			return nil, toResponse(questions, answers, m.formID)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *responsesModel) View() string {
	content :=
		styles.Heading.Render("Responses") +
			"\n" +
			baseStyle.Render(m.table.View()) +
			"\n\n" +
			styles.Description.Render("Enter - View Response | Press ESC to go back")

	return styles.PlaceCenter(m.width, m.height, content)
}

func (m *responsesModel) GetResponses() {
	formWithResponses, err := m.formSvc.GetWithResponses(m.formID, m.user)
	if err != nil {
		m.responseError = err
	}

	m.Form = *formWithResponses

	var order []string
	var columns []table.Column
	for _, question := range m.Form.Questions {
		columns = append(columns, table.Column{
			Title: question.Text,
			Width: m.width / (len(m.Form.Questions) + 2),
		})
		order = append(order, question.ID.String())
	}

	// TODO: This is a hack, need to fix this
	var rows []table.Row
	for _, response := range m.Form.Responses {
		var row table.Row
		for _, id := range order {
			for _, answer := range response.Answers {
				if answer.QuestionID.String() == id {
					row = append(row, answer.Value)
					break
				}
			}
		}

		rows = append(rows, row)
	}

	m.table.SetRows([]table.Row{})
	m.table.SetColumns([]table.Column{})

	m.table.SetColumns(columns)
	m.table.SetRows(rows)
}
