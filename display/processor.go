package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var coreContainer = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#ffffff")).
	BorderTop(true).
	BorderBottom(true).
	BorderLeft(true).
	BorderRight(true).
	Padding(0, 1)

func (m *displayModel) setupProcessorView(msg tea.WindowSizeMsg, width, height int) {
	coreContainer = coreContainer.Width(width/2 - 4).Height(height/2 - 2)
}

func (m *displayModel) updateProcessorViewSize(msg tea.WindowSizeMsg, width, height int) {
	coreContainer = coreContainer.Width(width/2 - 4).Height(height/2 - 2)
}

func (m *displayModel) updateProcessorView(msg tea.Msg) tea.Cmd {
	repr := strings.Builder{}
	core1 := coreContainer.Render("core 1")
	core2 := coreContainer.Render("core 2")
	core3 := coreContainer.Render("core 3")
	core4 := coreContainer.Render("core 4")
	m.processorView = lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, core1, core2),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, core3, core4),
	)
	repr.WriteString("\n")
	return nil
}
