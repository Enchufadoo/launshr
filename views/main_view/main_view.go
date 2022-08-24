package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/navigation"
	"launshr/parser"
	"launshr/shortcuts"
	"launshr/views/add_node"
	"launshr/views/command_list"
	"launshr/views/edit_node"
	"launshr/views/help"
	"os"
)

type ViewIndex int

const (
	CommandListView ViewIndex = iota
	EditNodeView
	AddNodeView
	HelpView
)

type Model struct {
	state          ViewIndex
	currentModel   tea.Model
	nodes          *parser.CommandNode
	configFilePath string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) setView(index ViewIndex) {
	m.state = index
}

func (m *Model) SetCurrentView(index ViewIndex) {
	m.setView(index)

	switch index {
	case HelpView:
		m.currentModel = help.New()
	case CommandListView:
		m.currentModel = command_list.New(m.nodes)
	case EditNodeView:
		m.currentModel = edit_node.New()
	case AddNodeView:
		m.currentModel = add_node.New()
	}

}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case shortcuts.HelpShortcut:
			if m.state == HelpView {
				m.SetCurrentView(CommandListView)
			} else {
				m.SetCurrentView(HelpView)
			}

		}
	case navigation.SaveCommandMsg:
		m.saveToFile(msg.(navigation.SaveCommandMsg).Node)
		m.SetCurrentView(CommandListView)
	case navigation.NavigateEditNodeViewMsg:
		m.SetCurrentView(EditNodeView)
	case navigation.NavigateAddNodeViewMsg:
		m.SetCurrentView(AddNodeView)
	case navigation.NavigateCommandListViewMsg:
		m.SetCurrentView(CommandListView)
	}

	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

func New(nodes *parser.CommandNode, configFilePath string) Model {
	clModel := command_list.New(nodes)
	return Model{
		nodes:          nodes,
		currentModel:   clModel,
		configFilePath: configFilePath,
	}
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.currentModel.View(),
	)
}

// The view shouldn't handle this kind of things, TODO create global app instance
func (m Model) saveToFile(node *parser.CommandNode) {
	err := parser.SaveToFile(node, m.configFilePath)
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
