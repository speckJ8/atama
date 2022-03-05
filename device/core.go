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
	ifetcher InstructionFetcher
	dline    DataLine
	// adder    Adder
	// multiplier Multiplier
	// divider Divider
	ICache    Cache
	DCache    Cache
	Name      string
	Mode      CoreMode
	Status    CoreStatus
	Registers RegisterSet
}

func NewCore(proc *Processor, mem *Memory, name string) Core {
	c := Core{proc: proc, mem: mem}
	c.ICache = NewCache(16, 4, 2)
	c.DCache = NewCache(16, 4, 2)
	c.Name = name
	return c
}

func (c *Core) Run() {
	// start fetching instructions from Registers.IP and executing
}

func (c *Core) ICacheSize() uint {
	return c.ICache.Size()
}

func (c *Core) DCacheSize() uint {
	return c.DCache.Size()
}
