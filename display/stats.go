package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var statsContainer = container.Copy()

func (m *displayModel) setupStatsView(msg tea.WindowSizeMsg, width, height int) {
	statsContainer = statsContainer.Width(width).Height(height - 1)
}

func (m *displayModel) updateStatsViewSize(msg tea.WindowSizeMsg, width, height int) {
	statsContainer = statsContainer.Width(width).Height(height - 1)
}

func (m *displayModel) updateStatsView(msg tea.Msg) tea.Cmd {
	repr := strings.Builder{}
	m.statsView = renderWithTitle(&statsContainer, boldText.Render("Stats"), 5, repr.String())
	return nil
}
