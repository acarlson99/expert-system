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

		// prog, err = Parse(toks)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }

		var prog interface{}
		var err error
		if toks[0][0] == 'v' { // TODO: remove
			verbose = verbose != true
			continue
		} else if toks[0][0] == 'h' {
			for ii, fact := range f.f {
				fmt.Printf("%s => %c\n", fact.rule, ii+'A')
			}
			continue
		} else {
			prog, err = Parse(toks)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		switch t := prog.(type) {
		case []Rule:
			for _, r := range t {
				f.AddRule(byte(r.id), r.node)
			}
		case Query:
			ret, err := f.UserQuery(t)
			if err != nil {
				panic(err)
			}
			fmt.Println(ret)
		case Assign:
			// TODO: remove panics
			err := f.UserSet(t)
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("%T: %+v\n", prog, prog)
			panic("Bad return from Parse")
		}
	}
	if scanner.Err() != nil {
		// TODO: hadle error
	}
}
