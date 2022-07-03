package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"launshr/components/form_element"
)

type Button struct {
	form_element.Selectable
	Text         string
	OnPressEnter func() tea.Cmd
}

func (b *Button) SetOnPressEnter(onPress func() tea.Cmd) {
	b.OnPressEnter = onPress
}

func (b *Button) GetText() string {
	return b.Text
}

func (b *Button) Update(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "enter":
			if b.OnPressEnter != nil {
				return b.OnPressEnter()
			}
		}
	}

	return nil
}

func (b *Button) SetText(text string) {
	b.Text = text
}

func (b *Button) Render() string {
	var style = lipgloss.NewStyle().
		Bold(true)

	if b.Selected == true {
		style.Foreground(lipgloss.Color("#cec5eb")).Bold(true)
	} else {
		style.Foreground(lipgloss.Color("#8f8f8f"))
	}

	return style.Render(b.Text)
}
