package command_list

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

type ListModel struct {
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
func (m ListModel) GenerateNodeModel(node *parser.CommandNode) ListModel {
	return ListModel{
		textInput:        m.textInput,
		filteredChildren: m.filterNodes(m.textInput.Value(), node),
		selected:         make(map[int]struct{}),
		currentNode:      node,
		cursor:           0,
		textFilterEmpty:  m.textInput.Value() == "",
	}
}

func (m ListModel) filterNodes(textToFilter string, node *parser.CommandNode) []parser.CommandNode {
	var filteredChildren []parser.CommandNode

	textFilter := strings.Trim(strings.ToLower(textToFilter), " ")
	textFilterEmpty := textFilter == ""

	for _, v := range node.Nodes {
		if !textFilterEmpty {
			if !strings.Contains(strings.ToLower(v.Name), textFilter) {
				continue
			}
		}

		filteredChildren = append(filteredChildren, v)
	}

	return filteredChildren
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

func InitialModel(node *parser.CommandNode) ListModel {
	newModel := ListModel{}

	filledModel := newModel.GenerateNodeModel(node)
	filledModel.textInput = generateTextInput()
	return filledModel
}

func (m ListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
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
			return ListModel{runningCommand: true}, runCommand(selectedNode.Command, selectedNode.WorkingDirectory)
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
func (m ListModel) View() string {

	if m.runningCommand {
		return ""
	}

	s := selectedItemStyle.Render("Command Menu")

	s += "\n"
	s += style.Render("")
	s += "\n"
	s += m.textInput.View()
	s += "\n\n"

	choicesString := "No results found"
	description := ""
	listItems := ""

	for i, choice := range m.filteredChildren {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		listItems += fmt.Sprintf("%s %s\n", cursor, choice.Name)
	}

	if len(m.filteredChildren) > 0 {
		choicesString = listItems
		description = RenderDescription(m.filteredChildren[m.cursor])
	}

	s += m.renderColumns(choicesString, description)

	return s
}

func (m ListModel) renderColumns(itemList string, description string) string {
	nameColumnWidth := 30
	nameColumnStyle := lipgloss.NewStyle().
		Margin(1, 3, 0, 0).
		Width(nameColumnWidth)

	descriptionColumnWidth := 50
	descriptionColumnStyle := lipgloss.NewStyle().
		Margin(1, 3, 0, 0).
		Width(descriptionColumnWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		nameColumnStyle.Copy().Render(itemList),
		descriptionColumnStyle.Copy().Render(description),
	)
}
