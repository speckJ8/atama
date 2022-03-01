package display

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/speckJ8/atama/device"
)

const (
	memoryContainerWidth = 40
)

type displayModel struct {
	processor     device.Processor
	memory        device.Memory
	ready         bool
	memoryView    string
	processorView string
	statusView    string
	cmdsView      string
}

func (m *displayModel) Init() tea.Cmd {
	m.memory = device.NewMemory(1 << 16)
	m.processor = device.NewProcessor(&m.memory)
	// return textinput.Blink
	return nil
}

func (m *displayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.setupMemoryView(msg, 26, msg.Height-2)
			m.setupProcessorView(msg, msg.Width-26, msg.Height-2)
			m.setupStatusView(msg, msg.Width)
			m.setupCmdsView(msg, msg.Width)
			m.ready = true
		} else {
			m.updateMemoryViewSize(msg, 26, msg.Height-2)
			m.updateProcessorViewSize(msg, msg.Width-26, msg.Height-1)
			m.updateStatusViewSize(msg, msg.Width)
			m.updateCmdsViewSize(msg, msg.Width)
		}
	}

	mCmd := m.updateMemoryView(msg)
	pCmd := m.updateProcessorView(msg)
	sCmd := m.updateStatusView(msg)
	cCmd := m.updateCmdsView(msg)
	if mCmd != nil {
		cmds = append(cmds, mCmd)
	}
	if pCmd != nil {
		cmds = append(cmds, pCmd)
	}
	if sCmd != nil {
		cmds = append(cmds, sCmd)
	}
	if cCmd != nil {
		cmds = append(cmds, cCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *displayModel) View() string {
	s := lipgloss.JoinHorizontal(lipgloss.Top, m.memoryView, m.processorView) +
		"\n" + m.statusView + "\n" + m.cmdsView
	return s
}

func Start() {
	p := tea.NewProgram(&displayModel{}, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
