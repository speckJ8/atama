package device

type OpCode = uint16

const (
	Nop OpCode = iota
	Halt
	// move between registers
	MovRR
	// move from memory to register
	MovMR
	// move from register to memory
	MovRM
	// move from immediate into register
	MovIR
	Add
	Mul
	Div
)

const (
	HaltMne = "halt"
	MovMne  = "mov"
	AddMne  = "add"
	MulMne  = "mul"
	DivMne  = "div"
)

var NopInstr = Instruction{OpCode: Nop}

type Instruction struct {
	OpCode   OpCode
	Operand1 Word
	Operand2 Word
	Operand3 Word
}
