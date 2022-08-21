package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
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

func (m *Model) showView(index ViewIndex, msg tea.Msg) {
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
				m.showView(CommandListView, msg)
			} else {
				m.showView(HelpView, msg)
			}

		}
	case navigation.SaveCommandMsg:
		m.saveToFile(msg.(navigation.SaveCommandMsg).Node)
		m.showView(CommandListView, msg)
	case navigation.NavigateEditNodeViewMsg:
		m.showView(EditNodeView, msg)
	case navigation.NavigateAddNodeViewMsg:
		m.showView(AddNodeView, msg)
	case navigation.NavigateCommandListViewMsg:
		m.showView(CommandListView, msg)
	}

	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

func InitialModel(nodes *parser.CommandNode, configFilePath string) Model {
	clModel := command_list.New(nodes)
	return Model{
		nodes:          nodes,
		currentModel:   clModel,
		configFilePath: configFilePath,
	}
}

func (m Model) View() string {
	return m.currentModel.View()
}

// The view shouldn't handle this kind of things, TODO create global app instance
func (m Model) saveToFile(node *parser.CommandNode) {
	err := parser.SaveToFile(node, m.configFilePath)
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
