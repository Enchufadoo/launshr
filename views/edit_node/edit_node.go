package edit_node

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/components/button"
	"launshr/components/form_element"
	"launshr/components/input"
	"launshr/navigation"
	"launshr/parser"
	"launshr/views/header"
)

const (
	NameInput = iota
	CommandInput
	WorkingDirectoryInput
	SaveButton
	CancelButton
)

type SaveNodeDataMsg struct{}
type JumpToNextItem struct{}

type Model struct {
	cursor          int
	selectedElement int
	listOfElements  *map[int]form_element.FormElement
	node            *parser.CommandNode
	header          header.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case SaveNodeDataMsg:
		return m.saveNodeData()
	case JumpToNextItem:
		m.moveCursorDown()
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+enter":
			return m.saveNodeData()
		case "esc":
			return m, navigation.EventNavigateCommandList()
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			m.moveCursorUp()
		case "down":
			m.moveCursorDown()
		}
	}

	for k, v := range *m.listOfElements {
		if k == m.cursor {
			cmd = v.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) View() string {

	return lipgloss.JoinVertical(lipgloss.Top,
		m.header.View(),
		"",
		(*m.listOfElements)[NameInput].Render(),
		(*m.listOfElements)[CommandInput].Render(),
		(*m.listOfElements)[WorkingDirectoryInput].Render(),
		"",
		(*m.listOfElements)[SaveButton].Render()+"\t"+(*m.listOfElements)[CancelButton].Render(),
	)

}

func (m *Model) SaveData() {
	m.node.Name = (*m.listOfElements)[NameInput].GetText()
	m.node.Command = (*m.listOfElements)[CommandInput].GetText()
	m.node.WorkingDirectory = (*m.listOfElements)[WorkingDirectoryInput].GetText()
}

func New(msg navigation.NavigateEditNodeViewMsg) tea.Model {
	nameElement := input.NewTextInput("Name",
		"Something to describe the command",
		inputPressEnterHandler)
	commandElement := input.NewTextInput("Command",
		"Command to run",
		inputPressEnterHandler)
	wdElement := input.NewTextInput("Working directory",
		"Working directory\", \"Path, empty is launshr CWD",
		inputPressEnterHandler)

	saveElement := &button.Button{
		Text: "Save",
		OnPressEnter: func() tea.Cmd {
			return func() tea.Msg {
				return SaveNodeDataMsg{}
			}
		},
	}
	cancelElement := &button.Button{
		Text: "Cancel",
		OnPressEnter: func() tea.Cmd {
			return navigation.EventNavigateCommandList()
		},
	}

	var listOfElements = map[int]form_element.FormElement{}
	listOfElements[NameInput] = nameElement
	listOfElements[CommandInput] = commandElement
	listOfElements[WorkingDirectoryInput] = wdElement
	listOfElements[SaveButton] = saveElement
	listOfElements[CancelButton] = cancelElement

	m := Model{
		listOfElements: &listOfElements,
		cursor:         0,
	}

	m.header = header.New()
	m.header.SubHeaderText = "Edit the command data"
	m.node = msg.Node
	(*m.listOfElements)[NameInput].SetText(msg.Node.Name)
	(*m.listOfElements)[CommandInput].SetText(msg.Node.Command)
	(*m.listOfElements)[WorkingDirectoryInput].SetText(msg.Node.WorkingDirectory)
	(*m.listOfElements)[NameInput].SetSelected(true)

	return m

}

func (m *Model) moveCursor(cursorDiff int) {
	m.cursor += cursorDiff

	for k, v := range *m.listOfElements {
		if k == m.cursor {
			v.SetSelected(true)
		} else {
			v.SetSelected(false)
		}
	}
}

func (m *Model) moveCursorDown() {
	if m.cursor < len(*m.listOfElements)-1 {
		m.moveCursor(1)
	}
}

func (m *Model) moveCursorUp() {
	if m.cursor > 0 {
		m.moveCursor(-1)
	}
}

func (m *Model) saveNodeData() (Model, tea.Cmd) {
	m.SaveData()
	return *m, navigation.EventSaveCommand(m.node)
}

func inputPressEnterHandler() tea.Cmd {
	return func() tea.Msg {
		return JumpToNextItem{}
	}
}
