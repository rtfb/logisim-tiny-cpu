package asm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rtfb/sketchbook/logisim/isa2/isa"
	"github.com/rtfb/sketchbook/logisim/isa2/parser"
)

func Assemble(input []parser.Token) ([]byte, error) {
	var opcode *isa.Opcode
	var out []byte
	for _, t := range input {
		if opcode == nil {
			op, ok := isa.ByName(t.Text)
			if !ok {
				return nil, fmt.Errorf("%s: unknown opcode '%s'", t.Position, t.Text)
			}
			opcode = &op
			if opcode.Param == isa.ParamIgnored {
				out = append(out, opcode.Emit(0))
				opcode = nil
			}
			continue
		}
		if opcode.Param != isa.ParamIsImmediate {
			reg, ok := isa.RegByName(t.Text)
			if !ok {
				return nil, fmt.Errorf("%s: bad register '%s' for opcode %s", t.Position, t.Text, opcode.Mnemonic)
			}
			out = append(out, opcode.Emit(reg.Code))
			opcode = nil
			continue
		}
		var immediate byte
		if strings.HasPrefix(t.Text, "0x") || strings.HasPrefix(t.Text, "0X") {
			imm, err := strconv.ParseInt(t.Text[2:], 16, 8)
			if err != nil {
				return nil, fmt.Errorf("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, opcode.Mnemonic, err)
			}
			immediate = byte(imm)
		} else if strings.HasPrefix(t.Text, "0") {
			imm, err := strconv.ParseInt(t.Text, 8, 8)
			if err != nil {
				return nil, fmt.Errorf("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, opcode.Mnemonic, err)
			}
			immediate = byte(imm)
		} else {
			imm, err := strconv.ParseInt(t.Text, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("%s: bad immediate argument '%s' for opcode %s: %v", t.Position, t.Text, opcode.Mnemonic, err)
			}
			immediate = byte(imm)
		}
		var limit byte
		opcode, limit = xformLoadImmediate(opcode, immediate)
		if immediate > limit {
			return nil, fmt.Errorf("%s: immediate argument '%s' for opcode %s is too large, must be <=7", t.Position, t.Text, opcode.Mnemonic)
		}
		out = append(out, opcode.Emit(immediate))
		opcode = nil
	}
	return out, nil
}

func xformLoadImmediate(inop *isa.Opcode, imm byte) (*isa.Opcode, byte) {
	if inop.Mnemonic == "li" {
		if imm <= 7 {
			return inop, 7
		}
		newOp, _ := isa.ByName("li1")
		return &newOp, 15
	}
	if inop.Mnemonic == "li0" || inop.Mnemonic == "li1" {
		return inop, 15
	}
	return inop, 7
}
