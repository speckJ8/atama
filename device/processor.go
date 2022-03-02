package device

import (
	"fmt"
)

type InterruptHandler = func()

type Processor struct {
	interruptVector map[uint]InterruptHandler
	mem             *Memory
	Cores           []Core
}

func NewProcessor(mem *Memory, numCores int) Processor {
	p := Processor{
		Cores: make([]Core, numCores),
	}
	for c := range p.Cores {
		p.Cores[c] = NewCore(&p, mem, fmt.Sprintf("Core %d", c))
	}
	return p
}

func (p *Processor) Start() {
	for c := range p.Cores {
		p.Cores[c].Status = CoreHalted
	}
	if len(p.Cores) > 0 {
		p.Cores[0].Status = CoreRunning
		p.Cores[0].Registers.IP = ZeroQWord
		go p.Cores[0].Run()
	}
}

func (p *Processor) NumCores() int {
	return len(p.Cores)
}
