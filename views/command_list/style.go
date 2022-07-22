package command_list

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

type ViewStyle struct {
	selectedItem      lipgloss.Style
	horizontalDivider lipgloss.Style
	fullWidth         lipgloss.Style
	fullWidthCenter   lipgloss.Style
	title             lipgloss.Style
}

func NewViewStyle() *ViewStyle {

	v := new(ViewStyle)

	v.fullWidth = lipgloss.NewStyle()

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	v.fullWidthCenter = lipgloss.NewStyle().Align(lipgloss.Center)

	if physicalWidth > 0 {
		v.fullWidth = v.fullWidth.Width(physicalWidth).MaxWidth(physicalWidth)
		v.fullWidthCenter = v.fullWidthCenter.Inherit(v.fullWidth)
	}

	v.selectedItem = lipgloss.NewStyle()
	v.horizontalDivider = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).
		BorderTop(true).BorderBottom(true).Inherit(v.fullWidth)

	v.title = v.fullWidthCenter.Copy().Bold(true)

	return v

}
