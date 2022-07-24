package command_list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/navigation"
	"launshr/parser"
	"launshr/shortcuts"
	"os"
	"os/exec"
	"strings"
)

type Model struct {
	currentNode     *parser.CommandNode
	cursor          int
	selected        map[int]struct{}
	textInput       textinput.Model
	textFilterEmpty bool
	runningCommand  bool
	children        *[]*parser.CommandNode
	vs              *ViewStyle
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

func InitialModel(node *parser.CommandNode) tea.Model {
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
			if m.cursor < len(*m.children)-1 {
				m.cursor++
			}
		case shortcuts.EDIT_COMMAND_SHORTCUT:
			return m, navigation.EventNavigateEditNode((*m.children)[m.cursor])
		case "backspace":
			return m.GenerateNodeModel(m.currentNode), cmd
		case "enter":
			if m.cursor < len(*m.children) {
				selectedNode := (*m.children)[m.cursor]
				if selectedNode.IsParent() {
					m.textInput.SetValue("")
					return m.GenerateNodeModel(selectedNode), cmd
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
		os.Exit(1)
		return programFinishedMsg{err} // TODO improve the exiting process
	})
}

// View Main view, represents the list of items to run
func (m Model) View() string {

	if m.runningCommand {
		return ""
	}

	s := lipgloss.JoinVertical(lipgloss.Left,
		m.renderHeader(),
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

func (m *Model) renderHeader() string {
	appName := m.vs.title.Render("Launshr")

	viewHeader := lipgloss.JoinVertical(lipgloss.Center, appName)

	var subtitleComponents []string

	if m.currentNode.Config != nil {
		if m.currentNode.Config.Title != "" {
			subtitleComponents = append(subtitleComponents, m.currentNode.Config.Title)
		}
	}

	if m.currentNode.IsParent() {
		if m.currentNode.Name != "" {
			subtitleComponents = append(subtitleComponents, m.currentNode.Name)
		}
	}

	subtitle := strings.Join(subtitleComponents, " - ")

	viewHeader = lipgloss.JoinVertical(lipgloss.Center, viewHeader, subtitle)

	return viewHeader
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

		listItems += fmt.Sprintf("%s %s\n", cursor, choice.Name)
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
