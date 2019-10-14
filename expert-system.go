package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

var verbose = false

func usage() {
	fmt.Println("usage: ./expert-system [options] [file]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = func() {
		usage()
		os.Exit(0)
	}

	var loadFile bool
	flag.BoolVar(&loadFile, "f", false, "Evaluate file and load repl")
	flag.BoolVar(&verbose, "v", false, "Verbose")

	flag.Parse()

	args := flag.Args()

	prog := GetFacts()

	switch len(args) {
	case 1:
		// file passed as arg
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("error: could not open file `%s`\n", args[0])
			os.Exit(1)
		}
		defer file.Close()
		err = ParseFile(file, prog, false)
		if err != nil {
			panic(err)
		}
		// only eval file or load data into repl
		if !loadFile {
			os.Exit(0)
		}
	case 0:
		break
	default:
		usage()
		os.Exit(1)
	}

	// enter REPL
	err := ParseFile(os.Stdin, prog, true)
	if err != nil {
		panic(err)
	}
}

func ParseFile(file *os.File, prog *Facts, setPrompt bool) error {
	rcloser := NewRCloser(file)
	var conf readline.Config
	conf.Stdin = rcloser
	if setPrompt {
		conf.Prompt = "> "
	}
	rl, err := readline.NewEx(&conf)
	defer rl.Close()
	for line, err := rl.Readline(); err == nil; line, err = rl.Readline() {
		eval(prog, line)
	}
	return err
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
		prog.UserSet(t)
	case Query:
		if len(t) > 0 {
			ret := prog.UserQuery(t)
			if err != nil {
				fmt.Println(err)
				return
			}
			tmp := strings.Join(strings.Fields(src[1:]), "")
			for ii, truth := range ret {
				if ii != 0 {
					fmt.Printf("\n")
				}
				fmt.Printf("[%c] = %v", tmp[ii], truth)
			}
			fmt.Println()
		}
	case []Rule:
		for _, r := range t {
			prog.AddRule(byte(r.id), r.node)
		}
	case Exit:
		os.Exit(0)
	case List:
		for i, f := range prog.f {
			str := ""
			if f.rule != nil {
				str = fmt.Sprintf("; %s => %c", f.rule.String(), i+'A')
			}
			fmt.Printf("[%c]: %t%s\n", i+'A', f.truth, str)
		}
	case Help:
		fmt.Printf("TODO: print funs")
	// TODO: handle cut
	case Vis:
		graph := prog.ToGraphviz()
		ast, err := graph.WriteAst()
		if err != nil {
			fmt.Println("Error creating graphviz AST:", err)
		}
		fmt.Println(ast)
	default:
		fmt.Printf("i-error: unknown parse return (%T,%+v)\n", ret, ret)
		return
	}
}
