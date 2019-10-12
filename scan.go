package main

import (
	"strings"
)

func uncomment(txt string) string {
	idx := strings.IndexRune(txt, '#')
	switch idx {
	case -1:
		return txt[:]
	default:
		return txt[:idx]
	}
}

func Scan(src string) string {
	return strings.Join(strings.Fields(uncomment(src)), " ")
}
