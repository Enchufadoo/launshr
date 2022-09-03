package add_node

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/components/form_element"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/edit_node/node_data_form"
	"launshr/views/header"
)

type SaveNodeDataMsg struct{}
type JumpToNextItem struct{}

type AddCommandMsg struct {
	Node *parser.CommandNode
	Msg  node_data_form.SaveNodeMsg
}

func EventAddCommand(node *parser.CommandNode, msg node_data_form.SaveNodeMsg) func() tea.Msg {
	return func() tea.Msg {
		return AddCommandMsg{
			Node: node,
			Msg:  msg,
		}
	}
}

type Model struct {
	cursor          int
	selectedElement int
	listOfElements  *map[int]form_element.FormElement
	editNodeForm    node_data_form.Model
	header          header.Model
	node            *parser.CommandNode
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.editNodeForm, cmd = m.editNodeForm.Update(msg)

	switch msg.(type) {
	case node_data_form.SaveNodeMsg:
		return m, EventAddCommand(m.node, msg.(node_data_form.SaveNodeMsg))
	}

	return m, cmd
}

func (m Model) View() string {
	s := lipgloss.JoinVertical(lipgloss.Top,
		m.header.View(),
		"",
		m.editNodeForm.View(),
	)

	return s
}

func New(msg navigation.NavigateAddNodeViewMsg) tea.Model {
	addHeader := header.New()
	addHeader.SubHeaderText = "Add a new command"

	dataForm := node_data_form.New()
	dataForm.InitializeForm()

	return Model{
		node:         msg.Parent,
		header:       addHeader,
		editNodeForm: dataForm,
	}
}

func inputPressEnterHandler() tea.Cmd {
	return func() tea.Msg {
		return JumpToNextItem{}
	}
}
