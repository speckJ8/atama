package device

type CoreStatus = uint8
type CoreMode = uint8

const (
	CoreHalted CoreStatus = iota
	CoreRunning
)

const (
	CorePrivileged CoreMode = iota
	CoreUnprivileged
)

type Core struct {
	proc     *Processor
	mem      *Memory
	icache   Cache
	dcache   Cache
	ifetcher InstructionFetcher
	dline    DataLine
	// adder    Adder
	// multiplier Multiplier
	// divider Divider
	Name      string
	Mode      CoreMode
	Status    CoreStatus
	Registers RegisterSet
}

func NewCore(proc *Processor, mem *Memory, name string) Core {
	c := Core{proc: proc, mem: mem}
	c.icache = NewCache(64, 4, 2)
	c.dcache = NewCache(64, 4, 2)
	c.Name = name
	return c
}

func (c *Core) Run() {
	// start fetching instructions from Registers.IP and executing
}
