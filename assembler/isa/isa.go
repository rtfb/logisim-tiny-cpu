package isa

type ParamType int

const (
	ParamIsRegister  ParamType = 0
	ParamIsImmediate ParamType = 1
	ParamIgnored     ParamType = 2
)

type Opcode struct {
	Code     byte   // binary value of the opcode extracted from the instruction
	Mnemonic string // the string representation of the instruction
	Param    ParamType
}

var ISA []Opcode = []Opcode{
	Opcode{
		Code:     0x00,
		Mnemonic: "halt",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x01,
		Mnemonic: "li",
		Param:    ParamIsImmediate,
	},
	Opcode{
		Code:     0x02,
		Mnemonic: "ld",
	},
	Opcode{
		Code:     0x03,
		Mnemonic: "st",
	},
	Opcode{
		Code:     0x04,
		Mnemonic: "getpc",
	},
	Opcode{
		Code:     0x05,
		Mnemonic: "getst",
	},
	Opcode{
		Code:     0x06,
		Mnemonic: "setst",
	},
	Opcode{
		Code:     0x07,
		Mnemonic: "shli",
		Param:    ParamIsImmediate,
	},
	Opcode{
		Code:     0x08,
		Mnemonic: "shri",
		Param:    ParamIsImmediate,
	},
	Opcode{
		Code:     0x09,
		Mnemonic: "getacc",
	},
	Opcode{
		Code:     0x0a,
		Mnemonic: "setacc",
	},
	Opcode{
		Code:     0x0b,
		Mnemonic: "swacc",
	},
	Opcode{
		Code:     0x0c,
		Mnemonic: "or",
	},
	Opcode{
		Code:     0x0d,
		Mnemonic: "and",
	},
	Opcode{
		Code:     0x0e,
		Mnemonic: "xor",
	},
	Opcode{
		Code:     0x0f,
		Mnemonic: "add",
	},
	Opcode{
		Code:     0x10,
		Mnemonic: "sub",
	},
	Opcode{
		Code:     0x11,
		Mnemonic: "inc",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x12,
		Mnemonic: "UNK",
		Param:    ParamIsImmediate,
	},
	Opcode{
		Code:     0x13,
		Mnemonic: "UNK",
		Param:    ParamIsImmediate,
	},
	Opcode{
		Code:     0x14,
		Mnemonic: "jz",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x15,
		Mnemonic: "jnz",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x16,
		Mnemonic: "jo",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x17,
		Mnemonic: "jno",
		Param:    ParamIgnored,
	},
	Opcode{
		Code:     0x18,
		Mnemonic: "jmp",
		Param:    ParamIgnored,
	},
}

func Lookup(code byte) Opcode {
	if int(code) > len(ISA)-1 {
		return Opcode{
			Code:     code,
			Mnemonic: "UNK",
			Param:    ParamIsImmediate,
		}
	}
	return ISA[code]
}
