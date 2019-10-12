package main

import (
	"bufio"
	"fmt"
	"os"
)

var verbose = false

func main() {
	f := GetFacts()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		toks := Scan(scanner.Text())
		prog, err := Parse(toks)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch t := prog.(type) {
		case []Rule:
			f.AddRule(byte(t[0].id), t[0].node)
		case []byte:
			switch t[0] {
			case '=':
				// TODO: assign
			case '?':
				f.UserQuery(t)
			default:
				panic("ahhh")
			}
		default:
			panic("Bad return from Parse")
		}
	}
	if scanner.Err() != nil {
		// TODO: hadle error
	}
}
