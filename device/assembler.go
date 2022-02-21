package device

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func Assemble(program []byte) ([]byte, error) {
	var text []byte
	var data []byte
	textLabels := make(map[string]int)
	dataLabels := make(map[string]int)
	relocations := make(map[string][]int)
	inTextSection := true
	P := len(program)
	line := 1
	for i := 0; i < P; {
		if unicode.IsSpace(rune(program[i])) {
			if program[i] == '\n' {
				line++
			}
			i++
			continue
		} else if program[i] == '.' {
			i++
			dir, off := parseIdentifier(program[i:])
			i += off
			switch dir {
			case "text":
				inTextSection = true
			case "data":
				inTextSection = false
			case "byte":
				if inTextSection {
					packByte(&text)
				} else {
					packByte(&data)
				}
			case "word":
				if inTextSection {
					packWord(&text)
				} else {
					packWord(&data)
				}
			case "dword":
				if inTextSection {
					packDWord(&text)
				} else {
					packDWord(&data)
				}
			case "qword":
				if inTextSection {
					packQWord(&text)
				} else {
					packQWord(&data)
				}
			default:
				return nil, errors.New(
					fmt.Sprintf("line %d: invalid directive %s", line, dir),
				)
			}
			continue
		}
		mne, off := parseIdentifier(program[i:])
		i += off
		// check if this is a label definition
		if i < P && program[i] == ':' {
			if inTextSection {
				textLabels[mne] = len(text)
			} else {
				dataLabels[mne] = len(data)
			}
			i++
			continue
		} else if !inTextSection {
			return nil, errors.New(fmt.Sprintf(
				"line %d: instructions can only be defined in a text section",
				line))
		}
		var opCode OpCode
		switch mne {
		case HaltMne:
			opCode = Halt
			packInstruction0(&text, opCode)
		case MovMne:
			op1, op1Type, off := parseOperand(program[i:], len(text)+2, relocations)
			i += off
			op2, op2Type, off := parseOperand(program[i:], len(text)+4, relocations)
			i += off
			if op1Type == 0 && op2Type == 0 {
				opCode = MovRR
			} else if op1Type == 0 && (op2Type == 1 || op2Type == 3) {
				opCode = MovRM
			} else if (op1Type == 1 || op1Type == 3) && op2Type == 0 {
				opCode = MovMR
			} else if op1Type == 2 && op2Type == 0 {
				opCode = MovIR
			} else {
				return nil, errors.New(
					fmt.Sprintf("line %d: illegal move instruction", line),
				)
			}
			packInstruction2(&text, opCode, op1, op2)
		case AddMne, MulMne, DivMne:
			op1, op1Type, off := parseOperand(program[i:], len(text)+2, relocations)
			i += off
			op2, op2Type, off := parseOperand(program[i:], len(text)+4, relocations)
			i += off
			if op1Type != 0 || op2Type != 0 {
				return nil, errors.New(
					fmt.Sprintf("line %d: `%s` operands can only be registers",
						line, mne),
				)
			}
			if mne == AddMne {
				opCode = Add
			} else if mne == MulMne {
				opCode = Mul
			} else {
				opCode = Div
			}
			packInstruction2(&text, opCode, op1, op2)
		default:
			return nil, errors.New(
				fmt.Sprintf("line %d: unknown mnemonic `%s`", line, mne),
			)
		}
	}
	for r := range relocations {
		rel := relocations[r]
		var pos int
		if _, ok := textLabels[r]; ok {
			pos = textLabels[r]
		} else if _, ok := dataLabels[r]; ok {
			pos = dataLabels[r] + len(text)
		} else {
			return nil, errors.New(fmt.Sprintf("undefined label %s", r))
		}
		posWord := Word{byte(pos), byte(pos >> 8)}
		for _, p := range rel {
			text[p] = posWord[0]
			text[p+1] = posWord[1]
		}
	}
	binary := append(text, data...)
	return binary, nil
}

func parseIdentifier(program []byte) (string, int) {
	P := len(program)
	i := 0
	for ; i < P && isAlphanumeric(program[i]); i++ {
	}
	return string(program[:i]), i
}

func parseOperand(program []byte, pos int, relocations map[string][]int) (Word, int, int) {
	P := len(program)
	i := 0
	for ; i < P && unicode.IsSpace(rune(program[i])); i++ {
	}
	j := i
	if i == P {
		return ZeroWord, -1, 0
	}
	if program[j] == '%' {
		i++
		j++
		for ; j < P && isAlphanumeric(program[j]); j++ {
		}
		reg := string(program[i:j])
		switch reg {
		case R0Name:
			return R0Number, 0, j
		case R1Name:
			return R1Number, 0, j
		case R2Name:
			return R2Number, 0, j
		case R3Name:
			return R3Number, 0, j
		case R4Name:
			return R4Number, 0, j
		case R5Name:
			return R5Number, 0, j
		case R6Name:
			return R6Number, 0, j
		case R7Name:
			return R7Number, 0, j
		case SPName:
			return SPNumber, 0, j
		case IPName:
			return IPNumber, 0, j
		default:
			return ZeroWord, -1, 0
		}
	} else if program[j] == '$' {
		i++
		j++
		for ; j < P && unicode.IsDigit(rune(program[j])); j++ {
		}
		imm, _ := strconv.Atoi(string(program[i:j]))
		return Word{byte(imm), byte(imm >> 8)}, 2, j
	} else if program[j] == '(' {
		i++
		j++
		base, stride, off := 0, 1, 0
		for ; j < P && unicode.IsDigit(rune(program[j])); j++ {
		}
		base, _ = strconv.Atoi(string(program[i:j]))
		if j != P && program[j] == ',' {
			j++
			i = j
			for ; j < P && unicode.IsDigit(rune(program[j])); j++ {
			}
			off, _ = strconv.Atoi(string(program[i:j]))
			if j != P && program[j] == ',' {
				j++
				i = j
				for ; j < P && unicode.IsDigit(rune(program[j])); j++ {
				}
				stride, _ = strconv.Atoi(string(program[i:j]))
			}
		}
		if j == P || program[j] != ')' {
			return ZeroWord, -1, 0
		}
		addr := base + stride*off
		return Word{byte(addr), byte(addr >> 8)}, 1, j + 1
	} else if unicode.IsLetter(rune(program[j])) {
		label, off := parseIdentifier(program[j:])
		j += off
		if relocations[label] == nil {
			relocations[label] = []int{pos}
		} else {
			relocations[label] = append(relocations[label], pos)
		}
		return Word{0, 0}, 3, j
	}

	return ZeroWord, -1, 0
}

func packInstruction0(text *[]byte, opCode OpCode) {
	*text = append(*text, byte(opCode), byte(opCode>>8), 0, 0, 0, 0, 0, 0)
}

func packInstruction2(text *[]byte, opCode OpCode, op1, op2 Word) {
	*text = append(*text, byte(opCode), byte(opCode>>8), op1[0], op1[1], op2[0], op2[1], 0, 0)
}

func packByte(binary *[]byte) {
	*binary = append(*binary, 0)
}

func packWord(binary *[]byte) {
	*binary = append(*binary, 0, 0)
}

func packDWord(binary *[]byte) {
	*binary = append(*binary, 0, 0, 0, 0)
}

func packQWord(binary *[]byte) {
	*binary = append(*binary, 0, 0, 0, 0, 0, 0, 0, 0)
}

func isAlphanumeric(b byte) bool {
	return unicode.IsDigit(rune(b)) || unicode.IsLetter(rune(b))
}
