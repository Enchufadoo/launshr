package command_list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/navigation"
	"launshr/parser"
	"launshr/shortcuts"
	"launshr/utils"
	"launshr/views/header"
	"os"
	"os/exec"
	"strings"
)

type Model struct {
	currentNode     *parser.CommandNode
	cursor          int
	selected        map[int]struct{}
	textInput       textinput.Model
	header          header.Model
	textFilterEmpty bool
	runningCommand  bool
	children        *[]*parser.CommandNode
	vs              *ViewStyle
}

type programFinishedMsg struct{ err error }

// GenerateNodeModel Used as a convenience method to update the model data
func (m Model) GenerateNodeModel(node *parser.CommandNode) Model {
	return Model{
		header:          m.header,
		textInput:       m.textInput,
		selected:        make(map[int]struct{}),
		currentNode:     node,
		cursor:          0,
		textFilterEmpty: m.textInput.Value() == "",
		children:        m.filterNodes(m.textInput.Value(), node),
		vs:              NewViewStyle(),
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
	ti.Prompt = "  "
	ti.SetCursorMode(textinput.CursorHide)

	return ti
}

func New(node *parser.CommandNode) tea.Model {
	newModel := Model{}

	filledModel := newModel.GenerateNodeModel(node)
	filledModel.textInput = generateTextInput()
	filledModel.header = header.New()
	return filledModel
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textInput, cmd = m.textInput.Update(msg)
	m.header, cmd = m.header.Update(msg)

	switch msg := msg.(type) {
	case programFinishedMsg:
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
				return m.GenerateNodeModel(m.currentNode.Parent),
					header.EventUpdateHeader(m.renderHeader(m.currentNode.Parent))
			}

		case "down":
			if m.cursor < len(*m.children)-1 {
				m.cursor++
			}
		case shortcuts.EditCommandShortcut:
			return m, navigation.EventNavigateEditNode((*m.children)[m.cursor])
		case shortcuts.AddCommandShortcut:
			return m, navigation.EventNavigateAddNode(m.currentNode)
		case "backspace":
			return m.GenerateNodeModel(m.currentNode), cmd
		case "enter":
			if m.cursor < len(*m.children) {
				selectedNode := (*m.children)[m.cursor]
				if selectedNode.IsParent() {
					m.textInput.SetValue("")
					return m.GenerateNodeModel(selectedNode), header.EventUpdateHeader(m.renderHeader(selectedNode))
				}
				return Model{runningCommand: true}, runCommand(selectedNode.Command, selectedNode.WorkingDirectory)
			}
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
		if err != nil {
			utils.ExitError(utils.ErrorRunningTheCommand, err)
		}

		return programFinishedMsg{err}
	})
}

// View Main view, represents the list of items to run
func (m Model) View() string {

	if m.runningCommand {
		return ""
	}

	s := lipgloss.JoinVertical(lipgloss.Left,
		m.header.View(),
		m.vs.horizontalDivider.Render(m.textInput.View()),
		m.renderItemsList(),
	)

	return s
}

func (m *Model) renderColumns(itemList string, description string, numberItems int) string {

	dividerLines := 4
	dividerString := ""

	if numberItems > 0 {
		if numberItems > dividerLines {
			dividerLines = numberItems
		}
		dividerString = m.vs.separatorColumnStyle.Render(strings.Repeat("\n", dividerLines))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.vs.nameColumnStyle.Render(itemList),
		dividerString,
		m.vs.descriptionColumnStyle.MarginLeft(1).Render(description),
	)
}

func (m *Model) renderHeader(node *parser.CommandNode) string {
	var subtitleComponents []string

	if node.Config != nil {
		if node.Config.Title != "" {
			subtitleComponents = append(subtitleComponents, node.Config.Title)
		}
	}

	if node.IsParent() {
		if node.Name != "" {
			subtitleComponents = append(subtitleComponents, node.Name)
		}
	}

	subtitle := strings.Join(subtitleComponents, " - ")

	return subtitle
}

func (m *Model) renderItemsList() string {
	choicesString := "No results found"
	description := ""
	listItems := ""

	for i, choice := range *m.children {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		name := choice.Name

		if choice.IsParent() {
			name = "\u2630 " + name
		}

		listItems += fmt.Sprintf("%s %s\n", cursor, name)
	}

	if len(*m.children) > 0 {
		choicesString = strings.Trim(listItems, "\n")
		description = RenderDescription(*(*m.children)[m.cursor])
	}

	headers := lipgloss.JoinHorizontal(lipgloss.Top,
		m.vs.nameHeader.Render("  Name"),
		m.vs.descriptionHeader.Render("  Description"))

	return lipgloss.JoinVertical(lipgloss.Top, headers,
		m.renderColumns(choicesString, description, len(*m.children)))

}
