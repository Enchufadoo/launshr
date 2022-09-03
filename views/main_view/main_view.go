package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/encoder"
	"launshr/navigation"
	"launshr/parser"
	"launshr/shortcuts"
	"launshr/views/add_node"
	"launshr/views/command_list"
	"launshr/views/edit_node"
	"launshr/views/help"
)

type Model struct {
	currentModel   tea.Model
	nodes          *parser.CommandNode
	configFilePath string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case shortcuts.HelpShortcut:
			_, ok := m.currentModel.(help.Model)
			if ok {
				m.NavigateToCommandList()
			} else {
				m.currentModel = help.New()
			}

		}
	case edit_node.EditCommandMsg:
		encoder.EditAndSaveToFile(msg.(edit_node.EditCommandMsg), m.configFilePath)
		m.NavigateToCommandList()
	case add_node.AddCommandMsg:
		m.nodes = encoder.AddAndSaveToFile(msg.(add_node.AddCommandMsg), m.configFilePath)
		m.NavigateToCommandList()
	case navigation.NavigateEditNodeViewMsg:
		m.currentModel = edit_node.New(msg.(navigation.NavigateEditNodeViewMsg))
	case navigation.NavigateAddNodeViewMsg:
		m.currentModel = add_node.New(msg.(navigation.NavigateAddNodeViewMsg))
	case navigation.NavigateCommandListViewMsg:
		m.NavigateToCommandList()
	}

	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

func (m *Model) NavigateToCommandList() {
	m.currentModel = command_list.New(m.nodes)
}

func New(nodes *parser.CommandNode, configFilePath string) Model {

	newModel := Model{
		nodes:          nodes,
		configFilePath: configFilePath,
	}

	newModel.NavigateToCommandList()

	return newModel
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.currentModel.View(),
	)
}
