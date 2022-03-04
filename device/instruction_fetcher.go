package device

type InstructionFetcher struct {
	icache *Cache
	mem    *Memory
	tasks  chan instructionFetchTask
}

type instructionFetchTask struct {
	Address      uint
	InstrChannel chan Instruction
	NextTask     *instructionFetchTask
}

func (f *InstructionFetcher) Init(icache *Cache, mem *Memory) {
	f.icache = icache
	f.mem = mem
	go func() {
		for task := range f.tasks {
			f.executeTask(&task)
		}
	}()
}

func (f *InstructionFetcher) Fetch(address uint) chan Instruction {
	c := make(chan Instruction)
	task := instructionFetchTask{
		Address:      address,
		InstrChannel: c,
	}
	f.tasks <- task
	return c
}

func (f *InstructionFetcher) executeTask(task *instructionFetchTask) {
	var instruction Instruction
	instr, cstatus := f.icache.GetQuadWord(task.Address)
	if cstatus == CacheAccessUnaligned {
		ShowCacheError(task.Address, cstatus)
		goto returnNop
	} else if cstatus == CacheAccessMiss {
		mstatus := f.icache.Populate(task.Address, f.mem)
		if mstatus == MemoryAccessInvalidAddress {
			ShowMemoryError(task.Address, mstatus)
			goto returnNop
		}
		instr, cstatus = f.icache.GetQuadWord(task.Address)
	}
	if cstatus == CacheAccessUnaligned {
		ShowCacheError(task.Address, cstatus)
		goto returnNop
	} else if cstatus == CacheAccessMiss {
		ShowCacheError(task.Address, cstatus)
		goto returnNop
	}

	instruction.OpCode = OpCode(uint16(instr[7])<<8 + uint16(instr[6]))
	instruction.Operand1 = Word{instr[4], instr[5]}
	instruction.Operand2 = Word{instr[2], instr[3]}
	instruction.Operand3 = Word{instr[0], instr[1]}
	task.InstrChannel <- instruction
	return

returnNop:
	task.InstrChannel <- NopInstr
	return
}
