package display

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/speckJ8/atama/device"
)

var coreContainer = container.Copy()

func (m *displayModel) setupProcessorView(msg tea.WindowSizeMsg, width, height int) {
	coreContainer = coreContainer.Width(width/2 - 3).Height(height)
}

func (m *displayModel) updateProcessorViewSize(msg tea.WindowSizeMsg, width, height int) {
	coreContainer = coreContainer.Width(width/2 - 3).Height(height)
}

func (m *displayModel) updateProcessorView(msg tea.Msg) tea.Cmd {
	core1 := renderWithTitle(&coreContainer, boldText.Render(m.processor.Cores[0].Name),
		len(m.processor.Cores[0].Name), coreView(&m.processor.Cores[0]))
	core2 := renderWithTitle(&coreContainer, boldText.Render(m.processor.Cores[1].Name),
		len(m.processor.Cores[1].Name), coreView(&m.processor.Cores[1]))
	m.processorView = lipgloss.JoinHorizontal(lipgloss.Top, core1, core2)
	return nil
}

func coreView(core *device.Core) string {
	repr := strings.Builder{}
	repr.WriteString("\n")
	repr.WriteString(cacheView(&core.ICache, "ICache"))
	repr.WriteString("\n")
	repr.WriteString(cacheView(&core.DCache, "DCache"))
	return repr.String()
}

func cacheView(cache *device.Cache, name string) string {
	repr := strings.Builder{}
	repr.WriteString(lightText.Render(
		fmt.Sprintf("--- %s %s",
			name, strings.Repeat("-", max(coreContainer.GetWidth()-7-len(name), 0))),
	))
	repr.WriteString("\n")
	repr.WriteString(lightText.Render(
		fmt.Sprintf("Addr  Set  Block%s  V  RU",
			strings.Repeat(" ", max(coreContainer.GetWidth()-25, 0))),
	))
	repr.WriteString("\n")
	for b := range cache.Blocks {
		block := cache.Blocks[b]
		set := (uint(b) % cache.SetSize) % cache.SetCount
		valid := "\u2022"
		if block.Valid {
			valid = greenText.Render(valid)
		} else {
			valid = redText.Render(valid)
		}
		ru := ""
		if block.RecentlyUsed {
			ru = "\u2022"
		}
		blockStr := strings.Builder{}
		for i := range block.Data {
			blockStr.WriteString(fmt.Sprintf("%02x ", block.Data[i]))
		}
		pad := strings.Repeat(" ",
			max(coreContainer.GetWidth()-20-3*int(cache.BlockSize), 0))
		repr.WriteString(
			fmt.Sprintf("%04x  %3d  %s%s  %s   %s\n",
				block.Address, set, blockStr.String(), pad, valid, ru),
		)
	}
	return repr.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
