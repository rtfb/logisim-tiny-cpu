package asm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rtfb/logisim-tiny-cpu/isa"
	"github.com/rtfb/logisim-tiny-cpu/parser"
)

func Assemble(input []parser.Token) ([]byte, error) {
	intermediate, labels, err := firstPass(input)
	if err != nil {
		return nil, err
	}
	return secondPass(intermediate, labels)
}

func firstPass(input []parser.Token) ([]intermOp, labelMap, error) {
	if len(input) == 0 {
		return nil, nil, nil
	}
	e := func(format string, a ...any) ([]intermOp, labelMap, error) {
		return nil, nil, fmt.Errorf(format, a...)
	}
	out := make([]intermOp, 255)
	seenLabels := make(labelMap)
	opIndex := 0
	addr := 0
	for _, t := range input {
		opcode := &out[opIndex]
		opcode.addr = addr
		if strings.HasSuffix(t.Text, ":") {
			seenLabels[strings.TrimRight(t.Text, ":")] = addr
			continue
		}
		if opcode.op.Empty() {
			op, ok := isa.ByName(t.Text)
			if !ok {
				return e("%s: unknown opcode '%s'", t.Position, t.Text)
			}
			if addr > 255 {
				return e("program too large: addr=%d", addr)
			}
			opcode.op = op
			if op.Param == isa.ParamIgnored {
				opIndex++
				addr++
			}
			continue
		}
		opcode.param = t
		opIndex++
		advance := opcode.op.ExpansionWidth
		if advance == 0 {
			advance = 1
		}
		addr += int(advance)
	}
	return out[:opIndex], seenLabels, nil
}

func secondPass(prog []intermOp, labels labelMap) ([]byte, error) {
	if len(prog) == 0 {
		return nil, nil
	}
	e := func(format string, a ...any) ([]byte, error) {
		return nil, fmt.Errorf(format, a...)
	}
	out := make([]byte, 0, 255)
	for _, interm := range prog {
		op := interm.op
		t := interm.param
		switch op.Param {
		case isa.ParamIsRegister:
			reg, ok := isa.RegByName(t.Text)
			if !ok {
				return e("%s: bad register '%s' for opcode %s", t.Position, t.Text, op.Mnemonic)
			}
			op.Emit(reg.Code)
			out = append(out, op.Emit(reg.Code))
		case isa.ParamIsImmediate:
			op, immediate, err := processImmediate(op, t)
			if err != nil {
				return nil, err
			}
			out = append(out, op.Emit(immediate))
		case isa.ParamIgnored:
			out = append(out, op.Emit(0))
		case isa.ParamIsLabel:
			jumpAddr, ok := labels[t.Text]
			if !ok {
				return e("%s: unknown label '%s' near opcode %s", t.Position, t.Text, op.Mnemonic)
			}
			if jumpAddr > 255 {
				return e("jumpAddr too large") // should never happen
			}
			out = emitJump(op, jumpAddr, out)
		}
	}
	return out, nil
}

func emitJump(op isa.Opcode, jumpAddr int, out []byte) []byte {
	highest4bits := byte((jumpAddr & 0xf0) >> 4)
	middleBit := byte(jumpAddr & 0x08)
	low3Bits := byte(jumpAddr & 0x07)
	liTmp, _ := isa.ByName("li")
	li, _ := xformLoadImmediate(liTmp, highest4bits)
	out = append(out, li.Emit(highest4bits))
	sjfName, sjfParam := decodeJumpFlavor(op.Mnemonic)
	sjf, _ := isa.ByName(sjfName)
	out = append(out, sjf.Emit(sjfParam))
	var jmp isa.Opcode
	if middleBit == 0 {
		jmp, _ = isa.ByName("jmplo")
	} else {
		jmp, _ = isa.ByName("jmphi")
	}
	return append(out, jmp.Emit(low3Bits))
}

// decodeJumpFlavor takes a jump pseudoinstruction mnemonic and returns an
// appropriate sjf or sjfn mnemonic to use, as well as an immediate
// representing the status flags to use that should go with it.
func decodeJumpFlavor(op string) (string, byte) {
	switch op {
	case "jz":
		return "sjf", 1
	case "jnz":
		return "sjfn", 1
	case "jo":
		return "sjf", 2
	case "jno":
		return "sjfn", 2
	default:
		return "sjf", 0
	}
}

func processImmediate(op isa.Opcode, t parser.Token) (isa.Opcode, byte, error) {
	e := func(format string, a ...any) (isa.Opcode, byte, error) {
		return isa.Opcode{}, 0, fmt.Errorf(format, a...)
	}
	var immediate byte
	if strings.HasPrefix(t.Text, "0x") || strings.HasPrefix(t.Text, "0X") {
		imm, err := strconv.ParseInt(t.Text[2:], 16, 8)
		if err != nil {
			return e("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, op.Mnemonic, err)
		}
		immediate = byte(imm)
	} else if strings.HasPrefix(t.Text, "0") {
		imm, err := strconv.ParseInt(t.Text, 8, 8)
		if err != nil {
			return e("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, op.Mnemonic, err)
		}
		immediate = byte(imm)
	} else {
		imm, err := strconv.ParseInt(t.Text, 10, 8)
		if err != nil {
			return e("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, op.Mnemonic, err)
		}
		immediate = byte(imm)
	}
	var limit byte
	op, limit = xformLoadImmediate(op, immediate)
	if immediate > limit {
		return e("%s: immediate argument '%s' for opcode %s is too large, must be <=7", t.Position, t.Text, op.Mnemonic)
	}
	return op, immediate, nil
}

func xformLoadImmediate(inop isa.Opcode, imm byte) (isa.Opcode, byte) {
	if inop.Mnemonic == "li" {
		if imm <= 7 {
			return inop, 7
		}
		newOp, _ := isa.ByName("li1")
		return newOp, 15
	}
	if inop.Mnemonic == "li0" || inop.Mnemonic == "li1" {
		return inop, 15
	}
	return inop, 7
}
