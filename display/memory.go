package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/speckJ8/atama/device"
)

var memoryContainer = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#ffffff")).
	BorderTop(true).
	BorderBottom(true).
	BorderLeft(true).
	BorderRight(true).
	Padding(0, 1)

var memoryViewport viewport.Model

func (m *DisplayModel) setupMemoryView(msg tea.WindowSizeMsg, width, height int) {
	memoryViewport = viewport.New(width, height-2)
	memoryContainer = memoryContainer.Width(width).Height(height - 2)
	// memoryViewport.HighPerformanceRendering = false
}
func (m *DisplayModel) updateMemoryViewSize(msg tea.WindowSizeMsg, width, height int) {
	memoryViewport.Height = height - 2
	memoryContainer = memoryContainer.Width(width).Height(height - 2)
}

func (m *DisplayModel) updateMemoryView(msg tea.Msg) tea.Cmd {
	repr := strings.Builder{}
	repr.WriteString("Memory\n")
	nLines := m.memory.Size / device.QWordSize
	for l := uint(0); l < nLines && l < 100; l++ {
		block := m.memory.GetQWordDirectly(l * device.QWordSize)
		repr.WriteString(fmt.Sprintf("%04x   %s\n", l, formatQuadWord(block)))
	}
	memoryViewport.SetContent(repr.String())
	m.memoryView = memoryContainer.Render(memoryViewport.View())
	return nil
}

func formatQuadWord(block device.QWord) string {
	// XXX: the block size is assumed to be 8 bytes here in order to avoid
	//      looping over block and perform multiple Sprintf's
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x%02x%02x",
		block[0], block[1], block[2], block[3],
		block[4], block[5], block[6], block[7])
}
