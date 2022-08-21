package add_node

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/components/form_element"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/edit_node/node_data_form"
)

type SaveNodeDataMsg struct{}
type JumpToNextItem struct{}

type Model struct {
	cursor          int
	selectedElement int
	listOfElements  *map[int]form_element.FormElement
	editNodeForm    tea.Model
	node            *parser.CommandNode
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.editNodeForm, cmd = m.editNodeForm.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	s := lipgloss.JoinVertical(lipgloss.Left,
		m.editNodeForm.View(),
	)

	return s
}

func (m *Model) SaveData() {

}

func New() tea.Model {
	return Model{
		editNodeForm: node_data_form.New(),
	}
}

func (m *Model) saveNodeData() (Model, tea.Cmd) {
	m.SaveData()
	return *m, navigation.EventSaveCommand(m.node)
}

func inputPressEnterHandler() tea.Cmd {
	return func() tea.Msg {
		return JumpToNextItem{}
	}
}
