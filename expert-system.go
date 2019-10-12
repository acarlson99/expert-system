package main

import (
	"bufio"
	"fmt"
	"os"
)

var verbose = false

func main() {
	prog := GetFacts()
	// read from file
	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("error: could not read from `%s`", os.Args[1])
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			eval(prog, scanner.Text())
		}
		if scanner.Err() != nil {
			fmt.Printf("error: %s\n", scanner.Err())
		}
	}
	// enter REPL
	scanner := bufio.NewScanner(os.Stdin)
	cnt := 0
	for readline(fmt.Sprintf("@%d: ", cnt), scanner) {
		cnt += 1
		eval(prog, scanner.Text())
		/* TODO: make work with new return
		if toks[0][0] == 'v' { // TODO: remove
			verbose = verbose != true
			continue
		} else if toks[0][0] == 'h' {
			for ii, fact := range f.f {
				fmt.Printf("%s => %c\n", fact.rule, ii+'A')
			}
			continue
		} else {
		*/

	}
	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

func readline(prompt string, scanner *bufio.Scanner) bool {
	fmt.Print(prompt)
	return scanner.Scan()
}

func eval(prog *Facts, src string) {
	ret, err := Parse(Scan(src))
	if err != nil {
		fmt.Println(err)
		return
	}
	switch t := ret.(type) {
	case nil:
		return
	case Assign:
		if len(t) == 0 {
			return
		}
		if prog.UserSet(t) != nil {
			fmt.Println(err)
			return
		}
	case Query:
		if len(t) == 0 {
			return
		}
		ret, err := prog.UserQuery(t)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ret)
	case []Rule:
		for _, r := range t {
			prog.AddRule(byte(r.id), r.node)
		}
	default:
		fmt.Printf("i-error: unknown parse return (%T,%+v)\n", ret, ret)
		return
	}
}
