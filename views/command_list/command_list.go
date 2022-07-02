package command_list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/navigation"
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
	currentNode     *parser.CommandNode
	cursor          int
	selected        map[int]struct{}
	textInput       textinput.Model
	textFilterEmpty bool
	runningCommand  bool
	children        *[]*parser.CommandNode
}

type programFinishedMsg struct{ err error }

// GenerateNodeModel Used as a convenience method to update the model data
func (m Model) GenerateNodeModel(node *parser.CommandNode) Model {
	return Model{
		textInput:       m.textInput,
		selected:        make(map[int]struct{}),
		currentNode:     node,
		cursor:          0,
		textFilterEmpty: m.textInput.Value() == "",
		children:        m.filterNodes(m.textInput.Value(), node),
	}
}

func (m Model) filterNodes(textToFilter string, node *parser.CommandNode) *[]*parser.CommandNode {
	var filteredChildren []*parser.CommandNode

	textFilter := strings.Trim(strings.ToLower(textToFilter), " ")
	textFilterEmpty := textFilter == ""

	for key, v := range node.Nodes {

		if !textFilterEmpty {
			if !strings.Contains(strings.ToLower(v.Name), textFilter) {
				continue
			}
		}

		filteredChildren = append(filteredChildren, &node.Nodes[key])
	}

	return &filteredChildren
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
			if m.cursor < len(*m.children)-1 {
				m.cursor++
			}
		case "ctrl+h":
			return m, navigation.EventNavigateEditNode((*m.children)[m.cursor])
		case "backspace":
			return m.GenerateNodeModel(m.currentNode), cmd
		case "enter":
			selectedNode := (*m.children)[m.cursor]
			if selectedNode.IsParent() {
				m.textInput.SetValue("")
				return m.GenerateNodeModel(selectedNode), cmd
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

	return tea.ExecProcess(c, func(err error) tea.Msg {
		os.Exit(1)
		return programFinishedMsg{err} // TODO improve the exiting process
	})
}

// View Main view, represents the list of items to run
func (m Model) View() string {

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

	for i, choice := range *m.children {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		listItems += fmt.Sprintf("%s %s\n", cursor, choice.Name)
	}

	if len(*m.children) > 0 {
		choicesString = listItems
		description = RenderDescription(*(*m.children)[m.cursor])
	}

	s += m.renderColumns(choicesString, description)

	return s
}

func (m Model) renderColumns(itemList string, description string) string {
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
