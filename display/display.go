package display

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/speckJ8/atama/device"
)

const (
	memoryContainerWidth = 32
)

type Display interface {
	Start()
	Refresh()
}

type displayModel struct {
	program       *tea.Program
	processor     *device.Processor
	memory        *device.Memory
	ready         bool
	memoryView    string
	processorView string
	statusView    string
	cmdsView      string
}

func (m *displayModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *displayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.setupMemoryView(msg, memoryContainerWidth, msg.Height-3)
			m.setupProcessorView(msg, msg.Width-memoryContainerWidth, msg.Height-4)
			m.setupStatusView(msg, msg.Width)
			m.setupCmdsView(msg, msg.Width)
			m.ready = true
		} else {
			m.updateMemoryViewSize(msg, memoryContainerWidth, msg.Height-3)
			m.updateProcessorViewSize(msg, msg.Width-memoryContainerWidth, msg.Height-4)
			m.updateStatusViewSize(msg, msg.Width)
			m.updateCmdsViewSize(msg, msg.Width)
		}
	}

	if m.ready {
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
	}
	return m, tea.Batch(cmds...)
}

func (m *displayModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	s := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.memoryView,
			m.processorView,
		),
		m.statusView,
		m.cmdsView,
	)
	return s
}

func NewDisplay(proc *device.Processor, mem *device.Memory) Display {
	return &displayModel{processor: proc, memory: mem}
}

func (m *displayModel) Refresh() {
	type nop struct{}
	if m.program != nil {
		m.program.Send(nop{})
	}
}

func (m *displayModel) Start() {
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	m.program = p
	if err := p.Start(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	os.Exit(0)
}
