package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/command_list"
	"launshr/views/edit_node"
	"os"
)

type ViewIndex int

const (
	CommandListView ViewIndex = iota
	EditNodeView
)

type Model struct {
	state            ViewIndex
	commandListModel command_list.Model
	editNodeModel    edit_node.Model
	nodes            *parser.CommandNode
	configFilePath   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case navigation.SaveCommandMsg:
		m.saveToFile(msg.(navigation.SaveCommandMsg).Node)
		m.state = CommandListView
	case navigation.NavigateEditNodeViewMsg:
		m.editNodeModel = edit_node.NewEditNodeModel()
		m.state = EditNodeView
	case navigation.NavigateCommandListViewMsg:
		m.state = CommandListView
	}

	switch m.state {
	case EditNodeView:
		m.editNodeModel, cmd = m.editNodeModel.Update(msg)
	case CommandListView:
		m.commandListModel, cmd = m.commandListModel.Update(msg)
	}

	return m, cmd
}

func InitialModel(node *parser.CommandNode, configFilePath string) Model {
	clModel := command_list.InitialModel(node)
	return Model{
		commandListModel: clModel,
		configFilePath:   configFilePath,
	}
}

func (m Model) View() string {
	switch m.state {
	case CommandListView:
		return m.commandListModel.View()
	case EditNodeView:
		return m.editNodeModel.View()
	default:
		return ""
	}
}

// The view shouldn't handle this kind of things, TODO create global app instance
func (m Model) saveToFile(node *parser.CommandNode) {
	err := parser.SaveToFile(node, m.configFilePath)
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
