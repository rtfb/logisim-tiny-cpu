package disasm

import (
	"fmt"
	"strings"
)

type paramType int

const (
	paramIsRegister  paramType = 0
	paramIsImmediate paramType = 1
	paramIgnored     paramType = 2
)

type opcode struct {
	code     byte   // binary value of the opcode extracted from the instruction
	mnemonic string // the string representation of the instruction
	param    paramType
}

var isa []opcode = []opcode{
	opcode{
		code:     0x00,
		mnemonic: "halt",
		param:    paramIgnored,
	},
	opcode{
		code:     0x01,
		mnemonic: "li",
		param:    paramIsImmediate,
	},
	opcode{
		code:     0x02,
		mnemonic: "ld",
	},
	opcode{
		code:     0x03,
		mnemonic: "st",
	},
	opcode{
		code:     0x04,
		mnemonic: "getpc",
	},
	opcode{
		code:     0x05,
		mnemonic: "getst",
	},
	opcode{
		code:     0x06,
		mnemonic: "setst",
	},
	opcode{
		code:     0x07,
		mnemonic: "shli",
		param:    paramIsImmediate,
	},
	opcode{
		code:     0x08,
		mnemonic: "shri",
		param:    paramIsImmediate,
	},
	opcode{
		code:     0x09,
		mnemonic: "getacc",
	},
	opcode{
		code:     0x0a,
		mnemonic: "setacc",
	},
	opcode{
		code:     0x0b,
		mnemonic: "swacc",
	},
	opcode{
		code:     0x0c,
		mnemonic: "or",
	},
	opcode{
		code:     0x0d,
		mnemonic: "and",
	},
	opcode{
		code:     0x0e,
		mnemonic: "xor",
	},
	opcode{
		code:     0x0f,
		mnemonic: "add",
	},
	opcode{
		code:     0x10,
		mnemonic: "sub",
	},
	opcode{
		code:     0x11,
		mnemonic: "inc",
	},
	opcode{
		code:     0x12,
		mnemonic: "UNK",
		param:    paramIsImmediate,
	},
	opcode{
		code:     0x13,
		mnemonic: "UNK",
		param:    paramIsImmediate,
	},
	opcode{
		code:     0x14,
		mnemonic: "jz",
		param:    paramIgnored,
	},
	opcode{
		code:     0x15,
		mnemonic: "jnz",
		param:    paramIgnored,
	},
	opcode{
		code:     0x16,
		mnemonic: "jo",
		param:    paramIgnored,
	},
	opcode{
		code:     0x17,
		mnemonic: "jno",
		param:    paramIgnored,
	},
	opcode{
		code:     0x18,
		mnemonic: "jmp",
		param:    paramIgnored,
	},
}

func lookup(code byte) opcode {
	if int(code) > len(isa)-1 {
		return opcode{
			code:     code,
			mnemonic: "UNK",
			param:    paramIsImmediate,
		}
	}
	return isa[code]
}

func fmtArg(opcode opcode, param byte) string {
	switch opcode.param {
	case paramIgnored:
		return ""
	case paramIsImmediate:
		return fmt.Sprintf(" 0x%d", param)
	case paramIsRegister:
		if param > 7 {
			return fmt.Sprintf(" 0x%d", param)
		}
		regs := []string{"r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7"}
		return " " + regs[param]
	default:
		return fmt.Sprintf(" 0x%d", param)
	}
}

// Do does the disassembling.
func Do(input []byte) string {
	var result strings.Builder
	for i, instruction := range input {
		opcode := lookup((instruction & 0xf8) >> 3)
		argument := instruction & 7
		result.WriteString(fmt.Sprintf("%3d: %s%s", i, opcode.mnemonic, fmtArg(opcode, argument)))
		if opcode.mnemonic == "UNK" {
			result.WriteString(fmt.Sprintf("  // 0b%b", instruction))
		}
		result.WriteString("\n")
	}
	return result.String()
}
