package device

import (
	"testing"
)

func TestAssembler(t *testing.T) {
	type AssemblerTest struct {
		asm    string
		binary []byte
		err    error
	}
	tests := []AssemblerTest{
		{
			asm: `
				mov $1 %r0
				mov $2 %r1
				add %r1 %r2
				halt
			`,
			binary: []byte{
				byte(MovIR), byte(MovIR >> 8), 1, 0, R0Number[0], R0Number[1], 0, 0,
				byte(MovIR), byte(MovIR >> 8), 2, 0, R1Number[0], R1Number[1], 0, 0,
				byte(Add), byte(Add >> 8), R1Number[0], R1Number[1],
				R2Number[0], R2Number[1], 0, 0,
				byte(Halt), byte(Halt >> 8), 0, 0, 0, 0, 0, 0,
			},
		},
		{
			asm: `
			.text
				mov a %r6
				mov b %r7
				mul %r6 %r7
			.data
			a:      .dword
			b:      .dword
			`,
			binary: []byte{
				byte(MovMR), byte(MovMR >> 8), 24, 0, R6Number[0], R6Number[1], 0, 0,
				byte(MovMR), byte(MovMR >> 8), 28, 0, R7Number[0], R7Number[1], 0, 0,
				byte(Mul), byte(Mul >> 8), R6Number[0], R6Number[1],
				R7Number[0], R7Number[1], 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}

	var binEq = func(a, b []byte) bool {
		if len(a) != len(b) {
			return false
		}
		for r := range a {
			if a[r] != b[r] {
				return false
			}
		}
		return true
	}

	for r := range tests {
		test := tests[r]
		binary, err := Assemble([]byte(test.asm))
		if err != test.err {
			if err == nil {
				t.Fatalf("test %d failed: expected error %s but obtained nil",
					r, test.err.Error())
			} else {
				t.Fatalf("test %d failed: %s", r, err.Error())
			}
		} else if !binEq(binary, test.binary) {
			t.Fatalf("test %d failed: obtained invalid binary", r)
		}
	}
}
