package display

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusContainer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#353533"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6124DF")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().Inherit(statusContainer).Padding(0, 1)
	infoStyle   = lipgloss.NewStyle().Padding(0, 1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#4273DB"))
)

func (m *displayModel) setupStatusView(msg tea.WindowSizeMsg, width int) {
	statusContainer = statusContainer.Width(width)
}

func (m *displayModel) updateStatusViewSize(msg tea.WindowSizeMsg, width int) {
	statusContainer = statusContainer.Width(width)
}

func (m *displayModel) updateStatusView(msg tea.Msg) tea.Cmd {
	w := lipgloss.Width
	title := titleStyle.Render("\u982d Atama")
	info := infoStyle.Render(
		fmt.Sprintf("Mem: %dKiB | 4 Cores", m.memory.Size/1024),
	)
	status := statusStyle.Width(statusContainer.GetWidth() - w(title) - w(info)).
		Render("This is the status of the application")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, title, status, info)
	m.statusView = statusContainer.Render(bar)
	return nil
}
