package device

type InterruptHandler = func()

type Processor struct {
	core1           Core
	core2           Core
	interruptVector map[uint]InterruptHandler
	mem             *Memory
}

func NewProcessor(mem *Memory) Processor {
	p := Processor{}
	p.core1 = NewCore(&p, mem)
	p.core2 = NewCore(&p, mem)
	return p
}

func (p *Processor) Start() {
	p.core1.Status = CoreRunning
	p.core2.Status = CoreHalted
	p.core1.Registers.IP = ZeroQWord
	go p.core1.Run()
}
