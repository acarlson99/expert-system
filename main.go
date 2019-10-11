package main

import (
	"bufio"
	"fmt"
	"os"
)

var verbose = false

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		toks := Scan(scanner.Text())
		prog, err := Parse(toks)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(prog)
	}
	if scanner.Err() != nil {
		// TODO: hadle error
	}
}
