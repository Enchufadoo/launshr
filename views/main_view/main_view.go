package main_view

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/parser"
	"launshr/views/command_list"
)

type state int

const (
	commandListView state = iota
)

type Model struct {
	state            state
	commandListModel command_list.ListModel
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case commandListView:
		var cmd tea.Cmd
		m.commandListModel, cmd = m.commandListModel.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func InitialModel(node *parser.CommandNode) Model {
	clModel := command_list.InitialModel(node)
	return Model{
		commandListModel: clModel,
	}
}

func (m Model) View() string {
	switch m.state {
	case commandListView:
		return m.commandListModel.View()
	default:
		return ""
	}
}
