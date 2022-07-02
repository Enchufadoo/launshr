package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/parser"
)

type NavigateEditNodeViewMsg struct {
	Node *parser.CommandNode
}

type NavigateCommandListViewMsg struct{}

func EventNavigateCommandList() func() tea.Msg {
	return func() tea.Msg {
		return NavigateCommandListViewMsg{}
	}
}

func EventNavigateEditNode(node *parser.CommandNode) func() tea.Msg {
	return func() tea.Msg {
		return NavigateEditNodeViewMsg{
			Node: node,
		}
	}
}
