package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rtfb/sketchbook/logisim/isa2/asm"
	"github.com/rtfb/sketchbook/logisim/isa2/disasm"
	"github.com/rtfb/sketchbook/logisim/isa2/parser"
	"github.com/rtfb/sketchbook/logisim/isa2/rom"
)

var disasmFlag bool

func init() {
	flag.BoolVar(&disasmFlag, "d", false, "disassemble")
}

func assemble(asmFilename string) {
	baseName := path.Base(asmFilename)
	input, err := ioutil.ReadFile(asmFilename)
	if err != nil {
		fmt.Printf("can't read %s: %v\n", asmFilename, err)
		return
	}
	tokens := parser.Tokenize(bytes.NewReader(input), baseName)
	output, err := asm.Assemble(tokens)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(rom.Dump(output))
	return
}

func disassemble(romFilename string) {
	rom, err := rom.Load(romFilename)
	if err != nil {
		panic(err)
	}
	noTrailingZeros := bytes.TrimRight(rom.Bytes, string([]byte{0}))
	assembly := disasm.Do(noTrailingZeros)
	fmt.Print(assembly)
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("Input file is required.\nEither a binary for disassemble, or an assembly file for assembling.\n")
		return
	}
	if disasmFlag {
		disassemble(flag.Args()[0])
	} else {
		assemble(flag.Args()[0])
	}
}
