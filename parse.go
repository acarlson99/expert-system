package main

import (
	"fmt"
	"strings"
)

type Exit struct{}
type List struct{}
type Vis struct {
	args []string
}
type Help struct{}

// FIXME:
// ABCDFEASFAFADSASFSADSADASFASFSAFADFAFASFSFDAFASFDSASFFSAADFDASADASFSAFDFDAAF|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||=>L
// ABCDEF=>K results in F=>K
// ?=A|B returns only 'B'

func Parse(src string) (interface{}, error) {
	if len(src) == 0 {
		return nil, nil
	} else if len(src) >= 2 && src[:2] == "?=" {
		in := strings.Join(strings.Fields(src[2:]), "")
		return parseOSQuery(in)
	} else if src[0] == '=' || src[0] == '?' {
		in := strings.Join(strings.Fields(src), "")
		return parseUnop(in)
	} else if idx := strings.Index(src, "=>"); idx != -1 {
		if idx == 0 {
			err := "left handside of rule `%s` is empty"
			return nil, fmt.Errorf("error: "+err, src)
		} else if len(src) == idx+2 {
			err := "right handside of rule `%s` is empty"
			return nil, fmt.Errorf("error: "+err, src)
		}
		lhs := strings.Join(strings.Fields(src[:idx]), "")
		rhs := strings.Join(strings.Fields(src[idx+2:]), "")
		return parseRule(src, lhs, rhs)
	} else if src == "v" || src == "verbose" {
		verbose = verbose != true
		fmt.Println("verbose =", verbose)
	} else if src == "x" || src == "exit" {
		return Exit{}, nil
	} else if src == "l" || src == "ls" || src == "list" {
		return List{}, nil
	} else if len(src) >= len("vis") && src[:3] == "vis" {
		in := strings.Fields(src)
		if len(in) > 1 {
			return Vis{args: in[1:]}, nil
		}
		return Vis{args: nil}, nil
	} else if src == "h" || src == "help" {
		return Help{}, nil
	} else {
		split := strings.Fields(src)
		if len(split) > 0 && (split[0] == "reset" || split[0] == "r") {
			return parseReset(split)
		}
		err := "unknown expression `%s`"
		return nil, fmt.Errorf("error: "+err, src)
	}
	return nil, nil
}

type Reset struct {
	args []byte
}

func parseReset(args []string) (Reset, error) {
	out := Reset{}
	if len(args) == 1 {
		for c := 'A'; c <= 'Z'; c += 1 {
			out.args = append(out.args, byte(c))
		}
	} else {
		for _, s := range args[1:] {
			if len(s) == 1 {
				c := byte(s[0])
				if c >= 'A' && c <= 'Z' {
					out.args = append(out.args, c)
				} else {
					err := "unknown literal `%s` in reset"
					return out, fmt.Errorf("error: "+err, s)
				}
			} else {
				err := "unknown literal `%s` in reset"
				return out, fmt.Errorf("error: "+err, s)
			}
		}
	}
	return out, nil
}

// @unop

func parseUnop(op string) (interface{}, error) {
	switch op[0] {
	case '=':
		return assign(op)
	case '?':
		return query(op)
	default:
		err := "unknown literal `%c` in expression `%s`"
		return nil, fmt.Errorf("error: "+err, op[0], op)
	}
}

type Assign []byte

func assign(src string) (Assign, error) {
	out := []byte{}
	if len(src) == 1 {
		return out, nil
	}
	for _, c := range src[1:] {
		if c >= 'A' && c <= 'Z' {
			out = append(out, byte(c))
		} else {
			err := "unknown literal `%c` in assignment `%s`"
			return out, fmt.Errorf("error: "+err, c, src)
		}
	}
	return out, nil
}

type Query []byte

func query(src string) (Query, error) {
	out := []byte{}
	if len(src) == 1 {
		return out, nil
	}
	for _, c := range src[1:] {
		if c >= 'A' && c <= 'Z' {
			out = append(out, byte(c))
		} else {
			err := "unknown literal `%c` in query `%s`"
			return out, fmt.Errorf("error: "+err, c, src)
		}
	}
	return out, nil
}

// @rule

type OSRule Rule

func parseOSQuery(src string) (OSRule, error) {
	out := OSRule{}
	if src == "" {
		return out, fmt.Errorf("error: empty rule")
	}
	if err := checkHs(src, src, "()!+|^", "right"); err != nil {
		return out, err
	}
	src1, err1 := toRPN(src)
	if err1 != nil {
		return out, err1
	}
	src2 := cleanNot(src1)
	src3 := simplifyHs(src2)
	out1, err2 := makeRule(src3, "A")
	if err2 != nil {
		return out, err2
	}
	if len(out1) != 1 {
		err := "invalid anonymous query `%s`"
		return out, fmt.Errorf("error: "+err, src)
	}
	return OSRule(out1[0]), nil
}

type Rule struct {
	id   rune
	node TreeNode
}

func parseRule(src string, lhs string, rhs string) ([]Rule, error) {
	out := []Rule{}
	if err := checkHs(src, lhs, "()!+|^", "left"); err != nil {
		return out, err
	}
	if err := checkHs(src, rhs, "()+", "right"); err != nil {
		return out, err
	}
	lhs1, err1 := toRPN(lhs)
	if err1 != nil {
		return out, err1
	}
	rhs1, err2 := toRPN(rhs)
	if err2 != nil {
		return out, err2
	}
	rhs1 = cleanNot(rhs1)
	lhs1 = cleanNot(lhs1)
	rhs2 := simplifyHs(rhs1)
	lhs2 := simplifyHs(lhs1)
	if err := checkRecDef(src, lhs2, rhs2); err != nil {
		return out, err
	}
	out1, err3 := makeRule(lhs1, rhs2)
	if err3 != nil {
		return out1, err3
	}
	return out1, nil
}

func cleanNot(src string) string {
	out := []byte{}
	for _, c := range src {
		if c != '#' {
			out = append(out, byte(c))
		}
	}
	return string(out)
}

func checkHs(src string, rule string, set string, hs string) error {
	if rule == "" {
		err := "empty %s handside in rule `%s`"
		return fmt.Errorf("error: "+err, hs, src)
	}
	for _, c := range rule {
		if c >= 'A' && c <= 'Z' || inSet(c, set) {
			continue
		} else {
			err := "unknown literal `%c` in %s handside of rule `%s`"
			return fmt.Errorf("error: "+err, c, hs, src)
		}
	}
	return nil
}

func inSet(c rune, s string) bool {
	for _, tmp := range s {
		if c == tmp {
			return true
		}
	}
	return false
}

func toRPN(infix string) (string, error) {
	prec := map[rune]int{
		'^': 1,
		'|': 2,
		'+': 4,
		'!': 8,
	}
	out_queue := []byte{}
	op_stack := []byte{}

	for _, c := range infix {
		right_assoc := c == '!'
		if c >= 'A' && c <= 'Z' {
			out_queue = append(out_queue, byte(c))
		} else if inSet(rune(c), "^|+!") {
			if c == '!' {
				out_queue = append(out_queue, '#')
			}
			p := prec[c]
			for len(op_stack) > 0 && ((prec[rune(op_stack[len(op_stack)-1])] > p) ||
				((prec[rune(op_stack[len(op_stack)-1])] == p) && !right_assoc)) &&
				(op_stack[len(op_stack)-1] != '(') {
				out_queue = append(out_queue, op_stack[len(op_stack)-1])
				op_stack = op_stack[:len(op_stack)-1]
			}
			op_stack = append(op_stack, byte(c))
		} else if c == '(' {
			op_stack = append(op_stack, byte(c))
		} else if c == ')' {
			for len(op_stack) > 0 && op_stack[len(op_stack)-1] != '(' {
				out_queue = append(out_queue, op_stack[len(op_stack)-1])
				op_stack = op_stack[:len(op_stack)-1]
			}
			if len(op_stack) == 0 {
				err := "mismatched parentheses in `%s`"
				return "", fmt.Errorf("error: "+err, infix)
			}
			if len(op_stack) > 0 && op_stack[len(op_stack)-1] == '(' {
				op_stack = op_stack[:len(op_stack)-1]
			}
		}
	}
	for len(op_stack) > 0 {
		c := op_stack[len(op_stack)-1]
		if inSet(rune(c), "()") {
			err := "mismatched parentheses in `%s`"
			return "", fmt.Errorf("error: "+err, infix)
		}
		out_queue = append(out_queue, byte(c))
		op_stack = op_stack[:len(op_stack)-1]
	}
	return string(out_queue), nil
}

func simplifyHs(rhs string) string {
	out := ""
	for _, c := range rhs {
		if c >= 'A' && c <= 'Z' {
			out += string(c)
		}
	}
	return out
}

func strrev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func makeRule(lhs string, rhs string) ([]Rule, error) {
	out := []Rule{}
	rrhs := strrev(rhs)
	for _, c := range rrhs {
		node, err := makeNode(lhs)
		if err != nil {
			return out, err
		}
		out = append(out, Rule{
			id:   c,
			node: node,
		})
	}
	return out, nil
}

func makeNode(lhs string) (TreeNode, error) {
	stack := []TreeNode{}
	var tree TreeNode
	var t1 TreeNode
	var t2 TreeNode
	_ = t1
	_ = t2
	for _, c := range lhs {
		if c >= 'A' && c <= 'Z' {
			tree = &Value{ch: byte(c)}
			stack = append(stack, tree)
		} else if c == '!' {
			tmp := &UnaryGate{
				gType: getType(c),
				next:  nil,
			}
			if len(stack) == 0 {
				return tree, fmt.Errorf("i-error: stack underflow in `makeNode`")
			}
			t2 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			tmp.next = t2
			tree = tmp
			stack = append(stack, tree)
		} else {
			tmp := &BinaryGate{
				gType: getType(c),
				left:  nil,
				right: nil,
			}
			if len(stack) == 0 {
				return tree, fmt.Errorf("i-error: stack underflow in `makeNode`")
			}
			t1 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				return tree, fmt.Errorf("i-error: stack underflow in `makeNode`")
			}
			t2 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			tmp.right = t1
			tmp.left = t2

			tree = tmp
			stack = append(stack, tree)
		}
	}
	if len(stack) == 0 {
		return tree, fmt.Errorf("i-error: stack underflow in `makeNode`")
	}
	tree = stack[len(stack)-1]
	return tree, nil
}

func getType(c rune) GType {
	switch c {
	case '!':
		return GateNot
	case '+':
		return GateAnd
	case '|':
		return GateOr
	case '^':
		return GateXor
	default:
		return GateNot
	}
}

func checkRecDef(src string, lhs string, rhs string) error {
	tmp := [26]int{0}
	for _, c := range lhs {
		if c >= 'A' && c <= 'Z' {
			tmp[int(c)-int('A')] += 1
		} else {
			err := "unknown literal `%c` in expression `%s`"
			return fmt.Errorf("error: "+err, c, src)
		}
	}
	for _, c := range rhs {
		if c >= 'A' && c <= 'Z' {
			if tmp[int(c)-int('A')] != 0 {
				err := "recursive definition of `%c` in expression `%s`"
				return fmt.Errorf("error: "+err, c, src)
			}
		} else {
			err := "unknown literal `%c` in expression `%s`"
			return fmt.Errorf("error: "+err, c, src)
		}
	}
	return nil
}
