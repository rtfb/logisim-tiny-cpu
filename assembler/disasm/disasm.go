package disasm

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/rtfb/logisim-tiny-cpu/isa"
)

func fmtArg(opcode isa.Opcode, param byte) string {
	switch opcode.Param {
	case isa.ParamIgnored:
		return ""
	case isa.ParamIsImmediate:
		return fmt.Sprintf(" 0x%d", param)
	case isa.ParamIsRegister:
		if param > 7 {
			return fmt.Sprintf(" 0x%d", param)
		}
		return " " + isa.Regs[param].Name
	default:
		return fmt.Sprintf(" 0x%d", param)
	}
}

// Do does the disassembling.
func Do(input []byte) string {
	var result strings.Builder
	w := tabwriter.NewWriter(&result, 3, 10, 1, '\t', 0)
	for i, instruction := range input {
		opcode := isa.Lookup((instruction & 0xf8) >> 3)
		argument := instruction & 7
		fmt.Fprintf(w, "%3d:\t%02x\t%s\t%s", i, instruction, opcode.Mnemonic, fmtArg(opcode, argument))
		if opcode.Mnemonic == "UNK" {
			fmt.Fprintf(w, "  // 0b%b", instruction)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return result.String()
}
