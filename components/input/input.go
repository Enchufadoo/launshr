package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"launshr/components/form_element"
)

type Input struct {
	form_element.Selectable
	input        textinput.Model
	prompt       string
	placeholder  string
	OnPressEnter func() tea.Cmd
}

func NewTextInput(prompt string, placeholder string, onPressEnter func() tea.Cmd) *Input {
	t := new(Input)
	t.prompt = prompt
	t.placeholder = placeholder
	t.OnPressEnter = onPressEnter
	t.generateTextInput()

	return t
}

func (i *Input) GetText() string {
	return i.input.Value()
}

func (i *Input) SetText(text string) {
	i.input.SetValue(text)
}

func (i *Input) Update(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "enter":
			if i.OnPressEnter != nil {
				return i.OnPressEnter()
			}

		}
	}

	input, cmd := i.input.Update(msg)
	i.input = input
	return cmd
}

func (i *Input) Render() string {
	retString := "  "

	if i.Selected {
		retString = "> "
	}

	return retString + i.input.View()
}

func (i *Input) SetSelected(selected bool) {

	if selected == i.Selected {
		return
	}

	i.Selected = selected

	if selected {
		i.input.SetCursorMode(textinput.CursorBlink)
		i.input.SetCursor(len(i.input.Value()))

	} else {
		i.input.SetCursorMode(textinput.CursorHide)
	}
}

func (i *Input) SetOnPressEnter(onPress func() tea.Cmd) {
	i.OnPressEnter = onPress
}

func (i *Input) generateTextInput() {
	ti := textinput.New()
	ti.Placeholder = i.placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Prompt = i.prompt + ": "
	ti.SetCursorMode(textinput.CursorHide)

	i.input = ti
}
