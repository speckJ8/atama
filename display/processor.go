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
	repr.WriteString(registersView(&core.Registers))
	repr.WriteString("\n")
	repr.WriteString(pipelineView(core))
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

func registersView(reg *device.RegisterSet) string {
	repr := strings.Builder{}
	repr.WriteString(lightText.Render(
		fmt.Sprintf("--- Registers %s", strings.Repeat("-", coreContainer.GetWidth()-16)),
	))
	repr.WriteString(
		fmt.Sprintf("%%r0=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R0[7], reg.R0[6], reg.R0[5], reg.R0[4],
			reg.R0[3], reg.R0[2], reg.R0[1], reg.R0[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r1=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R1[7], reg.R1[6], reg.R1[5], reg.R1[4],
			reg.R1[3], reg.R1[2], reg.R1[1], reg.R1[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r2=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R2[7], reg.R2[6], reg.R2[5], reg.R2[4],
			reg.R2[3], reg.R2[2], reg.R2[1], reg.R2[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r3=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R3[7], reg.R3[6], reg.R3[5], reg.R3[4],
			reg.R3[3], reg.R3[2], reg.R3[1], reg.R3[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r4=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R4[7], reg.R4[6], reg.R4[5], reg.R4[4],
			reg.R4[3], reg.R4[2], reg.R4[1], reg.R4[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r5=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R5[7], reg.R5[6], reg.R5[5], reg.R5[4],
			reg.R5[3], reg.R5[2], reg.R5[1], reg.R5[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r6=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R6[7], reg.R6[6], reg.R6[5], reg.R6[4],
			reg.R6[3], reg.R6[2], reg.R6[1], reg.R6[0]),
	)
	repr.WriteString(
		fmt.Sprintf("%%r7=%02x%02x%02x%02x%02x%02x%02x%02x  ",
			reg.R7[7], reg.R7[6], reg.R7[5], reg.R7[4],
			reg.R7[3], reg.R7[2], reg.R7[1], reg.R7[0]),
	)
	repr.WriteString(fmt.Sprintf("%%sp=%02x%02x  ", reg.SP[1], reg.SP[0]))
	repr.WriteString(fmt.Sprintf("%%ip=%02x%02x  ", reg.IP[1], reg.IP[0]))
	repr.WriteString("\n")
	return repr.String()
}

func pipelineView(core *device.Core) string {
	repr := strings.Builder{}
	repr.WriteString(lightText.Render(
		fmt.Sprintf("--- Pipeline %s", strings.Repeat("-", coreContainer.GetWidth()-15)),
	))
	repr.WriteString("\n")
	repr.WriteString(lightText.Render(
		"Stage         Instruction               Unit      Params                Cycles",
	))
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %-24s  %-8s  %-20s  %d",
			"Instr Fetch", "-", "IFetch", "0x1212", 4),
	)
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %-24s  %-8s  %-20s  %d",
			"Instr Decode",
			"-",
			"IDecode",
			"0x000200010006",
			2,
		),
	)
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %s%s  %-8s  %-20s  %d",
			"Data Fetch",
			boldText.Render("add %r1, %r2"),
			strings.Repeat(" ", 24-len("add %r0, %r1")),
			"DLine",
			"0x1212",
			4,
		),
	)
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %s%s  %-8s  %-20s  %d",
			"Execute",
			boldText.Render("add %r1, %r2"),
			strings.Repeat(" ", 24-len("add %r0, %r1")),
			"Adder 1",
			"100, 200",
			1,
		),
	)
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %s%s  %-8s  %-20s  %d",
			"",
			boldText.Render("div %r4, $2"),
			strings.Repeat(" ", 24-len("div %r4, $2")),
			"Divider",
			"8, 2",
			2,
		),
	)
	repr.WriteString("\n")
	repr.WriteString(
		fmt.Sprintf("%-12s  %s%s  %-8s  %-20s  %d",
			"Data Write",
			boldText.Render("mov %r3, $0x1212"),
			strings.Repeat(" ", 24-len("mov %r3, $0x1212")),
			"DLine",
			"0x500, 0x1212",
			4,
		),
	)
	repr.WriteString("\n\n")
	repr.WriteString(activeText.Render("IFetcher"))
	repr.WriteString("  ")
	repr.WriteString(activeText.Render("IDecoder"))
	repr.WriteString("  ")
	repr.WriteString(activeText.Render("Data Line"))
	repr.WriteString("  ")
	repr.WriteString(activeText.Render("Adder 1"))
	repr.WriteString("  ")
	repr.WriteString(inactiveText.Render("Adder 2"))
	repr.WriteString("  ")
	repr.WriteString(inactiveText.Render("Multiplier"))
	repr.WriteString("  ")
	repr.WriteString(activeText.Render("Divider"))
	repr.WriteString("\n")
	return repr.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
