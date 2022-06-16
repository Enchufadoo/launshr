package commandlist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/parser"
	"os"
	"os/exec"
	"strings"
)

var (
	selectedItemStyle = lipgloss.NewStyle()
	style             = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).Width(30)
)

type Model struct {
	currentNode      *parser.CommandNode
	filteredChildren []parser.CommandNode
	cursor           int
	selected         map[int]struct{}
	textInput        textinput.Model
	textFilterEmpty  bool
	runningCommand   bool
}

type programFinishedMsg struct{ err error }

// GenerateNodeModel Used as a convenience method to update the model data
func (m Model) GenerateNodeModel(node *parser.CommandNode) Model {
	var filteredChildren []parser.CommandNode

	textFilter := strings.Trim(strings.ToLower(m.textInput.Value()), " ")
	textFilterEmpty := textFilter == ""

	for _, v := range node.Nodes {
		if !textFilterEmpty {
			if !strings.Contains(strings.ToLower(v.Name), textFilter) {
				continue
			}
		}

		filteredChildren = append(filteredChildren, v)
	}

	return Model{
		textInput:        m.textInput,
		filteredChildren: filteredChildren,
		selected:         make(map[int]struct{}),
		currentNode:      node,
		cursor:           0,
		textFilterEmpty:  textFilterEmpty,
	}
}

func generateTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter commands"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = ""
	ti.SetCursorMode(textinput.CursorHide)

	return ti
}

func InitialModel(node *parser.CommandNode) Model {
	newModel := Model{}

	filledModel := newModel.GenerateNodeModel(node)
	filledModel.textInput = generateTextInput()
	return filledModel
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	switch msg := msg.(type) {
	case programFinishedMsg:
		os.Exit(0)
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "left":
			if m.currentNode.IsParent() {
				return m.GenerateNodeModel(m.currentNode.Parent), cmd
			}

		case "down":
			if m.cursor < len(m.filteredChildren)-1 {
				m.cursor++
			}
		case "backspace":
			if m.textFilterEmpty && m.currentNode.IsParent() {
				return m.GenerateNodeModel(m.currentNode.Parent), cmd
			}
			return m.GenerateNodeModel(m.currentNode), cmd
		case "enter":
			selectedNode := m.filteredChildren[m.cursor]
			if selectedNode.IsParent() {
				m.textInput.SetValue("")
				return m.GenerateNodeModel(&selectedNode), cmd
			}
			return Model{runningCommand: true}, runCommand(selectedNode.Command, selectedNode.WorkingDirectory)
		default:
			return m.GenerateNodeModel(m.currentNode), cmd
		}
	}

	return m, cmd
}

func runCommand(command string, workingDirectory string) tea.Cmd {
	if workingDirectory != "" {
		err := os.Chdir(workingDirectory)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	c := exec.Command("bash", "-c", command)

	wrapperCommand := tea.WrapExecCommand(c)

	return tea.Exec(wrapperCommand, func(err error) tea.Msg {
		os.Exit(1)
		return programFinishedMsg{err}
	})
}

// View Main view, represents the list of items to run
func (m Model) View() string {
	if m.runningCommand {
		return ""
	}

	columnWidth := 32

	columnsStyle := lipgloss.NewStyle().
		Margin(1, 3, 0, 0).
		Width(columnWidth)

	s := selectedItemStyle.Render("Command Menu")

	s += "\n"
	s += style.Render("")
	s += "\n"
	s += m.textInput.View()
	s += "\n\n"

	choicesString := ""

	for i, choice := range m.filteredChildren {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		choicesString += fmt.Sprintf("%s %s\n", cursor, choice.Name)
	}

	if choicesString == "" {
		choicesString = "No results found"
	}

	s += lipgloss.JoinHorizontal(
		lipgloss.Top,
		columnsStyle.Copy().Render(choicesString),
		columnsStyle.Copy().Render(RenderDescription(m.filteredChildren[m.cursor])),
	)

	return s
}
