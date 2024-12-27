package forms

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
	"github.com/devmegablaster/bashform/internal/styles"
)

type formsModel struct {
	width, height int
	Items         []models.Item
	list          list.Model
	session       ssh.Session

	isSubmitting  bool
	sizeError     bool
	submitError   error
	submitSuccess bool
	init          bool
	initTime      time.Time
}

func newFormsModel(items []models.Item, session ssh.Session) *formsModel {
	pty, _, _ := session.Pty()

	sizeErr := false
	if pty.Window.Width < 50 || pty.Window.Height < 30 {
		sizeErr = true
	}

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)

	l.SetSize(pty.Window.Width, 25)
	l.SetShowTitle(false)

	return &formsModel{
		Items:     items,
		list:      l,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		session:   session,
		init:      true,
		sizeError: sizeErr,
		initTime:  time.Now(),
	}
}

func (m *formsModel) Init() tea.Cmd {
	return nil
}

func (m *formsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			return m, toResponses(m.list.SelectedItem().(models.Item).ID)
		}
	}

	if m.init && time.Since(m.initTime) > 2*time.Second {
		m.init = false
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *formsModel) View() string {
	content := m.list.View()

	if m.init {
		return styles.PlaceCenter(m.width, m.height, constants.Logo)
	}

	box := styles.Box(m.width, content)

	return styles.PlaceCenterVertical(m.width, m.height, styles.Heading.Render("Your Forms"), box)
}
