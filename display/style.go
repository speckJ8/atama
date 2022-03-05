package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	text       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	boldText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	lightText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	greenText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	redText    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	activeText = lipgloss.NewStyle().Background(lipgloss.Color("#CCCCCC")).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1)
	inactiveText = lipgloss.NewStyle().Background(lipgloss.Color("#353533")).
			Foreground(lipgloss.Color("#CCCCCC")).
			Padding(0, 1)

	border = lipgloss.Border{}

	container = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ffffff")).
			BorderTop(false).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			Padding(0, 1)
)

func renderWithTitle(s *lipgloss.Style, title string, titlelen int, content string) string {
	header := fmt.Sprintf("┌─ %s %s┐", title, strings.Repeat("─", s.GetWidth()-titlelen-3))
	return lipgloss.JoinVertical(lipgloss.Top, header, s.Render(content))

}
