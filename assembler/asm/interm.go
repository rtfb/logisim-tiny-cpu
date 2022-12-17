package asm

import (
	"github.com/rtfb/sketchbook/logisim/isa2/isa"
	"github.com/rtfb/sketchbook/logisim/isa2/parser"
)

type intermOp struct {
	addr  int // address of this instruction
	op    isa.Opcode
	param parser.Token
}

func (op intermOp) equals(other intermOp) bool {
	return op.addr == other.addr && op.op == other.op && op.param.Text == other.param.Text
}

type labelMap map[string]int
