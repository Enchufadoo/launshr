package edit_node

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/edit_node/node_data_form"
	"launshr/views/header"
)

type SaveNodeDataMsg struct{}
type JumpToNextItem struct{}

type EditCommandMsg struct {
	Node *parser.CommandNode
	Msg  node_data_form.SaveNodeMsg
}

type Model struct {
	cursor          int
	selectedElement int
	node            *parser.CommandNode
	header          header.Model
	editNodeForm    node_data_form.Model
}

func EventEditCommand(node *parser.CommandNode, msg node_data_form.SaveNodeMsg) func() tea.Msg {
	return func() tea.Msg {
		return EditCommandMsg{
			Node: node,
			Msg:  msg,
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.editNodeForm, cmd = m.editNodeForm.Update(msg)

	switch msg.(type) {
	case node_data_form.SaveNodeMsg:
		return m, EventEditCommand(m.node, msg.(node_data_form.SaveNodeMsg))
	}

	return m, cmd
}

func (m Model) View() string {

	return lipgloss.JoinVertical(lipgloss.Top,
		m.header.View(),
		"",
		m.editNodeForm.View(),
	)

}

func New(msg navigation.NavigateEditNodeViewMsg) tea.Model {

	m := Model{
		cursor: 0,
	}

	dataForm := node_data_form.New()
	dataForm.InitializeForm()

	m.header = header.New()
	m.header.SubHeaderText = "Edit the command data"
	m.node = msg.Node

	m.editNodeForm = dataForm

	(*m.editNodeForm.ListOfElements)[node_data_form.NameInput].SetText(msg.Node.Name)
	(*m.editNodeForm.ListOfElements)[node_data_form.CommandInput].SetText(msg.Node.Command)
	(*m.editNodeForm.ListOfElements)[node_data_form.WorkingDirectoryInput].SetText(msg.Node.WorkingDirectory)
	(*m.editNodeForm.ListOfElements)[node_data_form.NameInput].SetSelected(true)

	return m

}
