package help

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/parser"
)

type Model struct {
}

// GenerateNodeModel Used as a convenience method to update the model data
func (m Model) GenerateNodeModel(node *parser.CommandNode) Model {
	return Model{}
}

func InitialModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View Main view, represents the list of items to run
func (m Model) View() string {
	result := ""
	subtleColor := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	bulletColor := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	subtitleStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)

	titleStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtleColor)

	shortcutBulletStyle := lipgloss.NewStyle().
		Foreground(bulletColor).
		PaddingRight(1)

	bullet := shortcutBulletStyle.SetString("*").String()

	result += titleStyle.Render("Launshr Help")

	result += subtitleStyle.Render("\nList of shortcuts")

	result += "\n" + bullet + "Show / Hide Help: F1 "
	result += "\n" + bullet + "Edit a command: CTRL + E "

	result += subtitleStyle.Render("\nProject Info: " +
		shortcutBulletStyle.Render("https://github.com/Enchufadoo/launshr"))

	return result
}
