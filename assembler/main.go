package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/rtfb/sketchbook/logisim/isa2/disasm"
	"github.com/rtfb/sketchbook/logisim/isa2/rom"
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
	romFilename := flag.Args()[0]
	rom, err := rom.Load(romFilename)
	if err != nil {
		panic(err)
	}
	noTrailingZeros := bytes.TrimRight(rom.Bytes, string([]byte{0}))
	assembly := disasm.Do(noTrailingZeros)
	fmt.Print(assembly)
}
