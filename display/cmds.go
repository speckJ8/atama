package display

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var cmdContainer = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5"))
var cmdInput textinput.Model

func (m *displayModel) setupCmdsView(msg tea.WindowSizeMsg, width int) {
	cmdInput = textinput.New()
	cmdInput.Placeholder = "Type a command..."
	cmdInput.Focus()
	cmdInput.Width = 20
	cmdContainer = cmdContainer.Width(width)
}

func (m *displayModel) updateCmdsViewSize(msg tea.WindowSizeMsg, width int) {
	cmdContainer = cmdContainer.Width(width)
}

func (m *displayModel) updateCmdsView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	cmdInput, cmd = cmdInput.Update(msg)
	m.cmdsView = cmdInput.View()
	return cmd
}
