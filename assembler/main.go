package main

import (
	"flag"
	"fmt"

	"github.com/rtfb/sketchbook/logisim/isa2/disasm"
)

var disasmFlag bool

func init() {
	flag.BoolVar(&disasmFlag, "d", false, "disassemble")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("Input file is required.\nEither a binary for disassemble, or an assembly file for assembling.\n")
		return
	}
	if !disasmFlag {
		fmt.Printf("Assembling is not implemented yet.\n")
		return
	}
	inputBytes := []byte{0x0f, 0x18, 0x0e, 0, 0xa0, 0xc0, 0x9b}
	assembly := disasm.Do(inputBytes)
	fmt.Print(assembly)
}
