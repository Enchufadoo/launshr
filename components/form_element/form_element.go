package form_element

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Selectable struct {
	Selected bool
}

func (s *Selectable) SetSelected(selected bool) {
	s.Selected = selected
}

func (s *Selectable) GetSelected() bool {
	return s.Selected
}

type FormElement interface {
	SetSelected(bool)
	GetSelected() bool
	Render() string
	Update(tea.Msg) tea.Cmd
	SetText(text string)
	GetText() string
}
