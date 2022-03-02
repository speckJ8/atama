package display

import (
	"fmt"
	"strings"

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

func (m *displayModel) setupMemoryView(msg tea.WindowSizeMsg, width, height int) {
	memoryContainer = memoryContainer.Width(width).Height(height - 2)
}
func (m *displayModel) updateMemoryViewSize(msg tea.WindowSizeMsg, width, height int) {
	memoryContainer = memoryContainer.Width(width).Height(height - 2)
}

func (m *displayModel) updateMemoryView(msg tea.Msg) tea.Cmd {
	repr := strings.Builder{}
	repr.WriteString("Memory")
	maxSize := m.memory.Size / device.QWordSize
	for l := uint(0); l < uint(memoryContainer.GetHeight()-1) && l < maxSize; l++ {
		block := m.memory.GetQWordDirectly(l * device.QWordSize)
		repr.WriteString(fmt.Sprintf("\n%04x   %s", l, formatQuadWord(block)))
	}
	m.memoryView = memoryContainer.Render(repr.String())
	return nil
}

func formatQuadWord(block device.QWord) string {
	// XXX: the block size is assumed to be 8 bytes here in order to avoid
	//      looping over block and perform multiple Sprintf's
	return fmt.Sprintf("%02x %02x %02x %02x %02x %02x %02x %02x",
		block[0], block[1], block[2], block[3],
		block[4], block[5], block[6], block[7])
}
