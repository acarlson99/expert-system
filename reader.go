package main

import (
	"fmt"
	"io"
	"strings"
)

const (
	LLParen   string = "("
	LRParen   string = ")"
	LNot      string = "!"
	LAnd      string = "+"
	LOr       string = "|"
	LXor      string = "^"
	LImply    string = "=>"
	LIfOnlyIf string = "<=>"
	LFact     string = "="
	LQuery    string = "?"
)

type Token struct {
	lexeme string
	lineno int
	column int
}

func Read(reader io.Reader) ([]Token, error) {
	const siz = 256
	buf := make([]byte, siz)
	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			break // TODO: return io error
		}
		str := strings.Split(string(buf[:siz]), "#")[0] // strip comments
		fmt.Println(str)
	}
	return nil, fmt.Errorf("TODO")
}
