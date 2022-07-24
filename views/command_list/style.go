package command_list

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

type ViewStyle struct {
	selectedItem           lipgloss.Style
	horizontalDivider      lipgloss.Style
	fullWidth              lipgloss.Style
	fullWidthCenter        lipgloss.Style
	title                  lipgloss.Style
	nameColumnStyle        lipgloss.Style
	descriptionColumnStyle lipgloss.Style
	separatorColumnStyle   lipgloss.Style
	nameHeader             lipgloss.Style
	descriptionHeader      lipgloss.Style
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

	nameColumnWidth := 30
	v.nameColumnStyle = lipgloss.NewStyle().
		Width(nameColumnWidth)

	descriptionColumnWidth := 50
	v.descriptionColumnStyle = lipgloss.NewStyle().
		Width(descriptionColumnWidth)

	v.separatorColumnStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).
		Border(lipgloss.NormalBorder(), false, false, false, true)

	listBorder := lipgloss.Border{Bottom: "―"}

	v.nameHeader = lipgloss.NewStyle().Width(30).
		BorderBottom(true).BorderStyle(listBorder).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#c7c7c7"})

	v.descriptionHeader = lipgloss.NewStyle().Width(physicalWidth - lipgloss.Width("  Description")).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).
		BorderBottom(true).BorderStyle(listBorder).
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#c7c7c7"})

	return v

}
