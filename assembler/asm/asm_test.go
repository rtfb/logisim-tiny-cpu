package asm

import (
	"errors"
	"strings"
	"testing"
	"text/scanner"

	"github.com/rtfb/sketchbook/logisim/isa2/isa"
	"github.com/rtfb/sketchbook/logisim/isa2/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirstPass(t *testing.T) {
	tests := map[string]struct {
		input      string
		want       []intermOp
		wantLabels labelMap
		wantErr    error
	}{
		"empty": {},
		"unknown_op": {
			input:   "foo",
			wantErr: errors.New("unknown_op:1:1: unknown opcode 'foo'"),
		},
		"label_and_op": {
			input: "label:\nli 1",
			want: []intermOp{
				{
					addr: 0,
					op:   isa.ISA[1],
					param: parser.Token{
						Text: "1",
					},
				},
			},
			wantLabels: labelMap{
				"label": 0,
			},
		},
		"several_ops": {
			input: `li 1
				loop:
					inc 5
					jmp loop
				after_loop:
					halt`,
			want: []intermOp{
				{
					addr: 0,
					op:   isa.ISA[1],
					param: parser.Token{
						Text: "1",
					},
				},
				{
					addr: 1,
					op:   isa.ISA[17],
					param: parser.Token{
						Text: "5",
					},
				},
				{
					addr: 2,
					op:   isa.ISA[24],
					param: parser.Token{
						Text: "loop",
					},
				},
				{
					addr: 5,
					op:   isa.ISA[0],
				},
			},
			wantLabels: labelMap{
				"loop":       1,
				"after_loop": 5,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tokens := parser.Tokenize(strings.NewReader(tt.input), name)
			got, gotLabels, err := firstPass(tokens)
			assert.Equal(t, tt.wantErr, err)
			// reset positions for easier assertion:
			for i := range got {
				got[i].param.Position = scanner.Position{}
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantLabels, gotLabels)
		})
	}
}

func TestSecondPass(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    []byte
		wantErr error
	}{
		"empty": {},
		"several_ops": {
			input: `li 1
				loop:
					inc 5
					jmp loop
				after_loop:
					halt`,
			want: []byte{
				0x09,     // li 1
				0x88 + 5, // inc 5
				0x08,     // li 0
				0xe0,     // sjf 0
				0xf1,     // jmplo 1
				0x00,     // halt
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tokens := parser.Tokenize(strings.NewReader(tt.input), name)
			firstPass, labels, err := firstPass(tokens)
			require.NoError(t, err)
			got, err := secondPass(firstPass, labels)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
