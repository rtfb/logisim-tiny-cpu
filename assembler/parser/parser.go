package parser

import (
	"io"
	"text/scanner"
)

type Token struct {
	Position scanner.Position
	Text     string
}

func Tokenize(input io.Reader, baseName string) []Token {
	var s scanner.Scanner
	s.Init(input)
	s.Filename = baseName
	var toks []Token
	var lastToken Token
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() == ":" {
			lastToken.Text += ":"
			s.Scan()
		}
		toks = append(toks, lastToken)
		lastToken = Token{
			Position: s.Position,
			Text:     s.TokenText(),
		}
	}
	return toks
}
