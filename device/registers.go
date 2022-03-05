package device

const (
	R0Name = "r0"
	R1Name = "r1"
	R2Name = "r2"
	R3Name = "r3"
	R4Name = "r4"
	R5Name = "r5"
	R6Name = "r6"
	R7Name = "r7"
	SPName = "sp"
	IPName = "ip"
)

var (
	R0Number = Word{0, 0}
	R1Number = Word{1, 0}
	R2Number = Word{2, 0}
	R3Number = Word{3, 0}
	R4Number = Word{4, 0}
	R5Number = Word{5, 0}
	R6Number = Word{6, 0}
	R7Number = Word{7, 0}
	SPNumber = Word{8, 0}
	IPNumber = Word{9, 0}
)

type RegisterSet struct {
	R0 QWord
	R1 QWord
	R2 QWord
	R3 QWord
	R4 QWord
	R5 QWord
	R6 QWord
	R7 QWord
	SP Word
	IP Word
}

func (s *RegisterSet) SetR0D(value DWord) {
	s.R0[0] = value[0]
	s.R0[1] = value[1]
	s.R0[2] = value[2]
	s.R0[3] = value[3]
}

func (s *RegisterSet) SetR0W(value Word) {
	s.R0[0] = value[0]
	s.R0[1] = value[1]
}

func (s *RegisterSet) SetR0B(value Byte) {
	s.R0[0] = value
}

func (s *RegisterSet) SetR1D(value DWord) {
	s.R1[0] = value[0]
	s.R1[1] = value[1]
	s.R1[2] = value[2]
	s.R1[3] = value[3]
}

func (s *RegisterSet) SetR1W(value Word) {
	s.R1[0] = value[0]
	s.R1[1] = value[1]
}

func (s *RegisterSet) SetR1B(value Byte) {
	s.R1[0] = value
}

func (s *RegisterSet) SetR2D(value DWord) {
	s.R2[0] = value[0]
	s.R2[1] = value[1]
	s.R2[2] = value[2]
	s.R2[3] = value[3]
}

func (s *RegisterSet) SetR2W(value Word) {
	s.R2[0] = value[0]
	s.R2[1] = value[1]
}

func (s *RegisterSet) SetR2B(value Byte) {
	s.R2[0] = value
}

func (s *RegisterSet) SetR3D(value DWord) {
	s.R3[0] = value[0]
	s.R3[1] = value[1]
	s.R3[2] = value[2]
	s.R3[3] = value[3]
}

func (s *RegisterSet) SetR3W(value Word) {
	s.R3[0] = value[0]
	s.R3[1] = value[1]
}

func (s *RegisterSet) SetR3B(value Byte) {
	s.R3[0] = value
}

func (s *RegisterSet) SetR4D(value DWord) {
	s.R4[0] = value[0]
	s.R4[1] = value[1]
	s.R4[2] = value[2]
	s.R4[3] = value[3]
}

func (s *RegisterSet) SetR4W(value Word) {
	s.R4[0] = value[0]
	s.R4[1] = value[1]
}

func (s *RegisterSet) SetR4B(value Byte) {
	s.R4[0] = value
}

func (s *RegisterSet) SetR5D(value DWord) {
	s.R5[0] = value[0]
	s.R5[1] = value[1]
	s.R5[2] = value[2]
	s.R5[3] = value[3]
}

func (s *RegisterSet) SetR5W(value Word) {
	s.R5[0] = value[0]
	s.R5[1] = value[1]
}

func (s *RegisterSet) SetR5B(value Byte) {
	s.R5[0] = value
}

func (s *RegisterSet) SetR6D(value DWord) {
	s.R6[0] = value[0]
	s.R6[1] = value[1]
	s.R6[2] = value[2]
	s.R6[3] = value[3]
}

func (s *RegisterSet) SetR6W(value Word) {
	s.R6[0] = value[0]
	s.R6[1] = value[1]
}

func (s *RegisterSet) SetR6B(value Byte) {
	s.R6[0] = value
}

func (s *RegisterSet) SetR7D(value DWord) {
	s.R7[0] = value[0]
	s.R7[1] = value[1]
	s.R7[2] = value[2]
	s.R7[3] = value[3]
}

func (s *RegisterSet) SetR7W(value Word) {
	s.R7[0] = value[0]
	s.R7[1] = value[1]
}

func (s *RegisterSet) SetR7B(value Byte) {
	s.R7[0] = value
}
