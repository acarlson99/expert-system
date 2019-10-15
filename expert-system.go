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
	case OSRule:
		fmt.Printf("%t\n", t.node.Evaluate())
	case []Rule:
		for _, r := range t {
			prog.AddRule(byte(r.id), r.node)
		}
	case Exit:
		os.Exit(0)
	case List:
		for i, f := range prog.f {
			str := ""
			if f.rule != nil || f.truth {
				if f.rule != nil {
					num := 1
					if f.truth {
						num = 2
					}
					str = fmt.Sprintf(";%*c%v => %c", num, ' ', f.rule, i+'A')
				}
				fmt.Printf("[%c]: %t%s\n", i+'A', f.truth, str)
			}
		}
	case Help:
		fmt.Printf(`=AB          Set A and B
?AB          Query A and B
?=(A | B)    Query expression
reset A B    Reset variable rules
list         List variables and rules
exit         Exit program
help         Display help
`)
	case Reset:
		for _, c := range t.args {
			if c >= 'A' && c <= 'Z' {
				prog.f[int(byte(c))-int('A')].rule = nil
			}
		}
	case Vis:
		fmt.Println("TODO: use args!!!")
		for _, s := range t.args {
			fmt.Println(s)
		}
		graph := prog.ToGraphviz()
		ast, err := graph.WriteAst()
		if err != nil {
			fmt.Println("Error creating graphviz AST:", err)
		}
		fmt.Println(ast)
		filename := "rules.dot"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error opening %s: %v\n", filename, err)
			return
		}
		defer file.Close()
		fmt.Printf("Dot form written to %s. Run `dot -Tpng %s -o rules.png` to generate png\n", filename, filename)
		file.Write([]byte(fmt.Sprintf("# dot -Tpng %s -o rules.png\n\n", filename)))
		file.Write([]byte(fmt.Sprintf("%v", ast)))
	default:
		fmt.Printf("i-error: unknown parse return (%T,%+v)\n", ret, ret)
		return
	}
}
