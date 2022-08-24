package help

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/views/header"
)

type Model struct {
	header header.Model
}

func New() Model {
	helpHeader := header.New()
	helpHeader.HeaderText = "Launshr Help"

	return Model{
		header: helpHeader,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View Main view, represents the list of items to run
func (m Model) View() string {

	bulletColor := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	subtitleStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)

	shortcutBulletStyle := lipgloss.NewStyle().
		Foreground(bulletColor).
		PaddingRight(1)

	bullet := shortcutBulletStyle.SetString("*").String()

	return lipgloss.JoinVertical(lipgloss.Top, m.header.View(),
		"List of shortcuts",
		"",
		bullet+"Show / Hide Help: F1 ",
		bullet+"Edit a command: CTRL + E ",
		bullet+"Add a command: CTRL + A ",
		subtitleStyle.Render("\nProject Info: "+
			shortcutBulletStyle.Render("https://github.com/Enchufadoo/launshr")),
	)

}
