package lexer

import (
	"io"
)

type Lexeme string

const (
	LParen   Lexeme = "("
	RParen   Lexeme = ")"
	Not      Lexeme = "!"
	And      Lexeme = "+"
	Or       Lexeme = "|"
	Xor      Lexeme = "^"
	Implies  Lexeme = "=>"
	IfOnlyIf Lexeme = "<=>"
)

type Token struct {
	lexeme Lexeme
	lineno int
	column int
}

func Read(Reader reader) ([]Token, err) {
	// TODO: implement :3
}
