package device

type DataType = uint8
type dataTaskType = uint8

const (
	DataQuadWord DataType = iota
	DataDoubleWord
	DataWord
	DataByte
)

const (
	DataRead dataTaskType = iota
	DataWrite
)

type DataLine struct {
	dcache *Cache
	mem    *Memory
	tasks  chan dataTask
}

type dataTask struct {
	Address     uint
	Type        DataType
	TaskType    dataTaskType
	Data        QWord
	DataChannel chan QWord
	NextTask    *dataTask
}

func (f *DataLine) Init(dcache *Cache, mem *Memory) {
	f.dcache = dcache
	f.mem = mem
	go func() {
		for task := range f.tasks {
			if task.TaskType == DataWrite {
				f.executeWriteTask(&task)
			} else {
				f.executeReadTask(&task)
			}
		}
	}()
}

func (f *DataLine) Read(address uint, t DataType) chan QWord {
	c := make(chan QWord)
	task := dataTask{
		Address:     address,
		Type:        t,
		TaskType:    DataRead,
		DataChannel: c,
	}
	f.tasks <- task
	return c
}

func (f *DataLine) Write(address uint, data QWord, t DataType) {
	task := dataTask{
		Address:     address,
		Type:        t,
		TaskType:    DataWrite,
		Data:        data,
		DataChannel: nil,
	}
	f.tasks <- task
}

func (f *DataLine) executeReadTask(task *dataTask) {
	var data QWord
	var cstatus CacheAccessStatus
	if task.Type == DataByte {
		var d Byte
		d, cstatus = f.dcache.GetByte(task.Address)
		data[0] = d
	} else if task.Type == DataWord {
		var d Word
		d, cstatus = f.dcache.GetWord(task.Address)
		data[0] = d[0]
		data[1] = d[1]
	} else if task.Type == DataDoubleWord {
		var d DWord
		d, cstatus = f.dcache.GetDoubleWord(task.Address)
		data[0] = d[0]
		data[1] = d[1]
		data[2] = d[2]
		data[3] = d[3]
	} else {
		data, cstatus = f.dcache.GetQuadWord(task.Address)
	}
	if cstatus == CacheAccessUnaligned {
		ShowCacheError(task.Address, cstatus)
		data = ZeroQWord
	} else if cstatus == CacheAccessMiss {
		mstatus := f.dcache.Populate(task.Address, f.mem)
		if mstatus == MemoryAccessInvalidAddress {
			ShowMemoryError(task.Address, mstatus)
			data = ZeroQWord
		}
		if task.Type == DataByte {
			var d Byte
			d, cstatus = f.dcache.GetByte(task.Address)
			data[0] = d
		} else if task.Type == DataWord {
			var d Word
			d, cstatus = f.dcache.GetWord(task.Address)
			data[0] = d[0]
			data[1] = d[1]
		} else if task.Type == DataDoubleWord {
			var d DWord
			d, cstatus = f.dcache.GetDoubleWord(task.Address)
			data[0] = d[0]
			data[1] = d[1]
			data[2] = d[2]
			data[3] = d[3]
		} else {
			data, cstatus = f.dcache.GetQuadWord(task.Address)
		}
	}
	if cstatus == CacheAccessUnaligned {
		ShowCacheError(task.Address, cstatus)
		data = ZeroQWord
	} else if cstatus == CacheAccessMiss {
		ShowCacheError(task.Address, cstatus)
		data = ZeroQWord
	}

	task.DataChannel <- data
	return
}

func (f *DataLine) executeWriteTask(task *dataTask) {
	var cstatus CacheAccessStatus
	if task.Type == DataByte {
		cstatus = f.dcache.SetByte(task.Address, task.Data[0])
	} else if task.Type == DataWord {
		cstatus = f.dcache.SetWord(task.Address, Word{task.Data[0], task.Data[1]})
	} else if task.Type == DataDoubleWord {
		cstatus = f.dcache.SetDoubleWord(task.Address,
			DWord{task.Data[0], task.Data[1], task.Data[2], task.Data[3]})
	} else {
		cstatus = f.dcache.SetQuadWord(task.Address, task.Data)
	}
	if cstatus == CacheAccessUnaligned {
		ShowCacheError(task.Address, cstatus)
	} else if cstatus == CacheAccessMiss {
		f.mem.Write(task.Address, task.Data[:])
	}
	return
}
