package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	HeaderText    string
	SubHeaderText string
	vs            *ViewStyle
}

type UpdateHeader struct {
	SubHeaderText string
}

func EventUpdateHeader(subHeaderText string) func() tea.Msg {
	return func() tea.Msg {

		return UpdateHeader{
			SubHeaderText: subHeaderText,
		}
	}
}

func New() Model {
	return Model{
		vs:         NewViewStyle(),
		HeaderText: "Launshr",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case UpdateHeader:

		m.SubHeaderText = msg.SubHeaderText
	}

	return m, cmd
}

func (m Model) View() string {

	appName := m.vs.title.Render(m.HeaderText)
	viewHeader := lipgloss.JoinVertical(lipgloss.Center, appName)

	viewHeader = lipgloss.JoinVertical(lipgloss.Center, viewHeader, m.SubHeaderText)
	return viewHeader
}
