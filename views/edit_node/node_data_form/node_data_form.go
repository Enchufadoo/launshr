package node_data_form

import (
	tea "github.com/charmbracelet/bubbletea"
	"launshr/components/button"
	"launshr/components/form_element"
	"launshr/components/input"
	"launshr/navigation"
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
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {

	case JumpToNextItem:
		m.moveCursorDown()
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+enter":

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

func inputPressEnterHandler() tea.Cmd {
	return func() tea.Msg {
		return JumpToNextItem{}
	}
}

func (m *Model) SaveData() {

}

func (m Model) View() string {
	viewString := ""
	viewString += (*m.listOfElements)[NameInput].Render() + "\n"
	viewString += (*m.listOfElements)[CommandInput].Render() + "\n"
	viewString += (*m.listOfElements)[WorkingDirectoryInput].Render() + "\n\n"

	viewString += (*m.listOfElements)[SaveButton].Render() + "\t" + (*m.listOfElements)[CancelButton].Render()
	return viewString
}

func New() tea.Model {
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

	nameElement.SetSelected(true)

	return Model{
		listOfElements: &listOfElements,
		cursor:         0,
	}
}
