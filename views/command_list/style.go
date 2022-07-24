package command_list

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"launshr/app_style"
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
		BorderForeground(lipgloss.AdaptiveColor(app_style.BackColor())).
		BorderTop(true).BorderBottom(true).Inherit(v.fullWidth)

	v.title = v.fullWidthCenter.Copy().Bold(true)

	nameColumnWidth := 30
	v.nameColumnStyle = lipgloss.NewStyle().
		Width(nameColumnWidth)

	descriptionColumnWidth := physicalWidth - nameColumnWidth
	v.descriptionColumnStyle = lipgloss.NewStyle().
		Width(descriptionColumnWidth)

	v.separatorColumnStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.AdaptiveColor(app_style.BackColor())).
		Border(lipgloss.NormalBorder(), false, false, false, true)

	listBorder := lipgloss.Border{Bottom: "―"}

	v.nameHeader = lipgloss.NewStyle().Width(nameColumnWidth).
		BorderBottom(true).BorderStyle(listBorder).
		BorderForeground(lipgloss.AdaptiveColor(app_style.BackColor())).
		Bold(true).
		Foreground(lipgloss.AdaptiveColor(app_style.SecondaryColor()))

	v.descriptionHeader = lipgloss.NewStyle().Width(descriptionColumnWidth).
		BorderForeground(lipgloss.AdaptiveColor(app_style.BackColor())).
		BorderBottom(true).BorderStyle(listBorder).
		Bold(true).
		Foreground(lipgloss.AdaptiveColor(app_style.SecondaryColor()))

	return v

}
