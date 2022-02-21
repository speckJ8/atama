package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/speckJ8/atama/device"
)

const (
	memoryContainerWidth = 40
)

type DisplayModel struct {
	processor     device.Processor
	memory        device.Memory
	ready         bool
	memoryView    string
	processorView string
}

func (m *DisplayModel) Init() tea.Cmd {
	m.memory = device.NewMemory(1 << 16)
	m.processor = device.NewProcessor(&m.memory)
	return nil
}

func (m *DisplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.setupMemoryView(msg, 26, msg.Height)
			m.setupProcessorView(msg, msg.Width-26, msg.Height)
			m.ready = true
		} else {
			m.updateMemoryViewSize(msg, 26, msg.Height)
			m.updateProcessorViewSize(msg, msg.Width-26, msg.Height)
		}
	}

	mCmd := m.updateMemoryView(msg)
	pCmds := m.updateProcessorView(msg)
	cmds = append(cmds, mCmd)
	cmds = append(cmds, pCmds...)
	return m, tea.Batch(cmds...)
}

func (m *DisplayModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.memoryView, m.processorView)
}
