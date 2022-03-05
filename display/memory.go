package display

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/speckJ8/atama/device"
)

var memoryContainer = container.Copy()

func (m *displayModel) setupMemoryView(msg tea.WindowSizeMsg, width, height int) {
	memoryContainer = memoryContainer.Width(width).Height(height - 1)
}

func (m *displayModel) updateMemoryViewSize(msg tea.WindowSizeMsg, width, height int) {
	memoryContainer = memoryContainer.Width(width).Height(height - 1)
}

func (m *displayModel) updateMemoryView(msg tea.Msg) tea.Cmd {
	repr := strings.Builder{}
	maxSize := m.memory.Size / device.QWordSize
	for l := uint(0); l < uint(memoryContainer.GetHeight()-1) && l < maxSize; l++ {
		block := m.memory.GetQWordDirectly(l * device.QWordSize)
		repr.WriteString(fmt.Sprintf("\n%04x   %s", l, formatQuadWord(block)))
	}
	m.memoryView = renderWithTitle(&memoryContainer, boldText.Render("Memory"), 6,
		repr.String())
	return nil
}

func formatQuadWord(block device.QWord) string {
	// XXX: the block size is assumed to be 8 bytes here in order to avoid
	//      looping over block and perform multiple Sprintf's
	return fmt.Sprintf("%02x %02x %02x %02x %02x %02x %02x %02x",
		block[0], block[1], block[2], block[3],
		block[4], block[5], block[6], block[7])
}
