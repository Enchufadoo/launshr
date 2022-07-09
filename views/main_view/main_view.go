package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/command_list"
	"launshr/views/edit_node"
	"launshr/views/help"
	"os"
)

type ViewIndex int

const (
	CommandListView ViewIndex = iota
	EditNodeView
	HelpView
)

type Model struct {
	state          ViewIndex
	currentModel   tea.Model
	nodes          *parser.CommandNode
	configFilePath string
	lastView       ViewIndex
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) setView(index ViewIndex) {
	m.lastView = m.state
	m.state = index
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "f2":
			if m.state == HelpView {
				m.currentModel = command_list.InitialModel(m.nodes)
				m.setView(CommandListView)
			} else {
				m.currentModel = help.InitialModel()
				m.setView(HelpView)
			}

		}
	case navigation.SaveCommandMsg:
		m.saveToFile(msg.(navigation.SaveCommandMsg).Node)
		m.setView(CommandListView)
	case navigation.NavigateEditNodeViewMsg:
		m.currentModel = edit_node.InitialModel()
		m.setView(EditNodeView)
	case navigation.NavigateCommandListViewMsg:
		m.currentModel = command_list.InitialModel(m.nodes)
		m.setView(CommandListView)
	}

	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

func InitialModel(nodes *parser.CommandNode, configFilePath string) Model {
	clModel := command_list.InitialModel(nodes)
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
