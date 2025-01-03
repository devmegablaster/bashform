package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	Error       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	Succes      = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	Heading     = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)
	Description = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))
	Logo        = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)
)

func TableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#16a34a")).
		Bold(true)

	return s
}

func Box(width int, content string) string {
	boxWidth := 90
	if width < 90 {
		boxWidth = width - 5
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#94a3b8")).
		Align(lipgloss.Center).
		Width(boxWidth).
		Padding(1, 2).
		Render(content)
}

func PlaceCenterVertical(width, height int, content ...string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			content...,
		),
	)
}

func PlaceCenter(width, height int, content string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
