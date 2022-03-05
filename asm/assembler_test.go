package asm

import (
	"testing"

	"github.com/speckJ8/atama/device"
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
				byte(device.MovIR), byte(device.MovIR >> 8), 1, 0,
				device.R0Number[0], device.R0Number[1], 0, 0,
				byte(device.MovIR), byte(device.MovIR >> 8), 2, 0,
				device.R1Number[0], device.R1Number[1], 0, 0,
				byte(device.Add), byte(device.Add >> 8),
				device.R1Number[0], device.R1Number[1],
				device.R2Number[0], device.R2Number[1], 0, 0,
				byte(device.Halt), byte(device.Halt >> 8), 0, 0, 0, 0, 0, 0,
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
				byte(device.MovMR), byte(device.MovMR >> 8), 24, 0,
				device.R6Number[0], device.R6Number[1], 0, 0,
				byte(device.MovMR), byte(device.MovMR >> 8), 28, 0,
				device.R7Number[0], device.R7Number[1], 0, 0,
				byte(device.Mul), byte(device.Mul >> 8),
				device.R6Number[0], device.R6Number[1],
				device.R7Number[0], device.R7Number[1], 0, 0,
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
