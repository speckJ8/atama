package device

type memoryTaskType = uint8
type MemoryAccessStatus = uint8

const (
	memoryTaskRead  memoryTaskType = iota
	memoryTaskWrite memoryTaskType = iota
)

const (
	MemoryAccessInvalidAddress MemoryAccessStatus = iota
	MemoryAccessWriteDone      MemoryAccessStatus = iota
	MemoryAccessReadDone       MemoryAccessStatus = iota
)

type Memory struct {
	Size  uint
	data  []byte
	tasks chan memoryTask
}

type memoryTask struct {
	Type          memoryTaskType
	Address       uint
	Data          []byte
	Size          uint
	StatusChannel chan MemoryAccessStatus
	DataChannel   chan []byte
	NextTask      *memoryTask
}

func NewMemory(size uint) Memory {
	m := Memory{Size: size}
	m.data = make([]byte, size)
	m.tasks = make(chan memoryTask)
	return m
}

func (m *Memory) Init() {
	go func() {
		for task := range m.tasks {
			if task.Address+task.Size >= uint(len(m.data)) {
				task.StatusChannel <- MemoryAccessInvalidAddress
			} else if task.Type == memoryTaskRead {
				data := m.data[task.Address : task.Address+task.Size]
				task.StatusChannel <- MemoryAccessReadDone
				task.DataChannel <- data
			} else if task.Type == memoryTaskWrite {
				copy(m.data[task.Address:task.Address+task.Size], task.Data)
			}
		}
	}()
}

func (m *Memory) Write(address uint, data []byte) chan MemoryAccessStatus {
	c := make(chan MemoryAccessStatus)
	task := memoryTask{
		Type:          memoryTaskWrite,
		Address:       address,
		Data:          data,
		Size:          uint(len(data)),
		StatusChannel: c,
	}
	m.tasks <- task
	return c
}

func (m *Memory) Read(address uint, size uint) (chan MemoryAccessStatus, chan []byte) {
	c1 := make(chan []byte)
	c2 := make(chan MemoryAccessStatus)
	task := memoryTask{
		Type:          memoryTaskRead,
		Address:       address,
		Size:          size,
		DataChannel:   c1,
		StatusChannel: c2,
	}
	m.tasks <- task
	return c2, c1
}

func (m *Memory) GetQWordDirectly(addr uint) QWord {
	return QWord{
		m.data[addr], m.data[addr+1], m.data[addr+2], m.data[addr+3],
		m.data[addr+4], m.data[addr+5], m.data[addr+6], m.data[addr+7],
	}
}
