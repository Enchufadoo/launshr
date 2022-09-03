package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/parser"
)

type NavigateAddNodeViewMsg struct {
	Parent *parser.CommandNode
}

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

func EventNavigateAddNode(parent *parser.CommandNode) func() tea.Msg {
	return func() tea.Msg {
		return NavigateAddNodeViewMsg{
			Parent: parent,
		}
	}
}
