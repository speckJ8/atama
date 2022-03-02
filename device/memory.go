package device

import (
	"sync"
)

type MemoryTaskType = uint8
type MemoryAccessStatus = uint8

const (
	MemoryTaskRead  MemoryTaskType = iota
	MemoryTaskWrite MemoryTaskType = iota
)

const (
	MemoryAccessInvalidAddress MemoryAccessStatus = iota
	MemoryAccessWriteDone      MemoryAccessStatus = iota
	MemoryAccessReadDone       MemoryAccessStatus = iota
)

type Memory struct {
	Size       uint
	data       []byte
	taskQMutex *sync.Mutex
	taskQHead  *MemoryTask
	taskQTail  *MemoryTask
}

type MemoryTask struct {
	Type          MemoryTaskType
	Address       uint
	Data          []byte
	Size          uint
	StatusChannel chan MemoryAccessStatus
	DataChannel   chan []byte
	NextTask      *MemoryTask
}

func NewMemory(size uint) Memory {
	m := Memory{Size: size}
	m.data = make([]byte, size)
	m.taskQMutex = &sync.Mutex{}
	return m
}

func (m *Memory) Init() {
	go func() {
		for {
			if task := m.popTask(); task != nil {
				if task.Address+task.Size >= uint(len(m.data)) {
					task.StatusChannel <- MemoryAccessInvalidAddress
				} else if task.Type == MemoryTaskRead {
					data := m.data[task.Address : task.Address+task.Size]
					task.StatusChannel <- MemoryAccessReadDone
					task.DataChannel <- data
				} else if task.Type == MemoryTaskWrite {
					copy(m.data[task.Address:task.Address+task.Size], task.Data)
				}
			}
		}
	}()
}

func (m *Memory) Write(address uint, data []byte) chan MemoryAccessStatus {
	c := make(chan MemoryAccessStatus)
	task := MemoryTask{
		Type:          MemoryTaskWrite,
		Address:       address,
		Data:          data,
		Size:          uint(len(data)),
		StatusChannel: c,
	}
	m.pushTask(&task)
	return c
}

func (m *Memory) Read(address uint, size uint) (chan MemoryAccessStatus, chan []byte) {
	c1 := make(chan []byte)
	c2 := make(chan MemoryAccessStatus)
	task := MemoryTask{
		Type:          MemoryTaskRead,
		Address:       address,
		Size:          size,
		DataChannel:   c1,
		StatusChannel: c2,
	}
	m.pushTask(&task)
	return c2, c1
}

func (m *Memory) pushTask(task *MemoryTask) {
	m.taskQMutex.Lock()
	defer m.taskQMutex.Unlock()
	if m.taskQTail != nil {
		m.taskQTail.NextTask = task
	} else {
		m.taskQHead = task
	}
	m.taskQTail = task
}

func (m *Memory) popTask() *MemoryTask {
	m.taskQMutex.Lock()
	defer m.taskQMutex.Unlock()
	if m.taskQHead == nil {
		return nil
	}
	task := m.taskQHead
	m.taskQHead = m.taskQHead.NextTask
	if m.taskQHead == nil {
		m.taskQTail = nil
	}
	return task
}

func (m *Memory) GetQWordDirectly(addr uint) QWord {
	return QWord{
		m.data[addr], m.data[addr+1], m.data[addr+2], m.data[addr+3],
		m.data[addr+4], m.data[addr+5], m.data[addr+6], m.data[addr+7],
	}
}
