package main

import (
	"fmt"
	"strings"
)

func Parse(toks []string) (interface{}, error) {
	switch len(toks) {
	case 0:
		return nil, nil
	case 1:
		return parseUnop(toks[0])
	default:
		rule := false
		for _, s := range toks {
			if s == "=>" {
				rule = true
				break
			}
		}
		if rule {
			return parseRule(toks)
		} else {
			estr := strings.Join(toks, " ")
			err := "unknown expression `%s`"
			return nil, fmt.Errorf("error: "+err, estr)
		}
	}
}

// @unop

func parseUnop(op string) ([]byte, error) {
	switch op[0] {
	case '=':
		return assign(op)
	case '?':
		return query(op)
	default:
		// TODO: handle builtins
		return nil, nil
	}
}

func assign(src string) ([]byte, error) {
	out := []byte{}
	out = append(out, src[0])
	if len(src) == 1 {
		err := "empty left handside in assignment"
		return out, fmt.Errorf("error: " + err)
	}
	for _, c := range src[1:] {
		fmt.Println(c)
		if c >= 'A' && c <= 'Z' {
			out = append(out, byte(c))
		} else {
			err := "unknown literal `%c` in assignment `%s`"
			return out, fmt.Errorf("error: "+err, c, src)
		}
	}
	return out, nil
}

func query(src string) ([]byte, error) {
	out := []byte{}
	out = append(out, src[0])
	if len(src) == 1 {
		err := "empty left handside in query"
		return out, fmt.Errorf("error: " + err)
	}
	for _, c := range src[1:] {
		fmt.Println(c)
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

type Rule struct {
	id   rune
	node TreeNode
}

func parseRule(toks []string) ([]Rule, error) {
	out := []Rule{}
	lhs, rhs := splitHs(toks)
	if err := checkHs(toks, lhs, "()!+|^", "left"); err != nil {
		return out, err
	}
	if err := checkHs(toks, rhs, "()+", "right"); err != nil {
		return out, err
	}
	lhs1, err1 := toRPN(toks, lhs, "left")
	if err1 != nil {
		return out, err1
	}
	rhs1, err2 := toRPN(toks, rhs, "left")
	if err2 != nil {
		return out, err2
	}
	rhs2 := simplifyRhs(rhs1)
	// TODO: check against recursive definitons?
	out = makeRule(lhs1, rhs2)
	return out, nil
}

func splitHs(toks []string) (string, string) {
	lhs := ""
	rhs := ""
	left := true
	for _, s := range toks {
		if s == "=>" {
			toks = toks[:1]
			left = false
		} else if left {
			lhs += s
			toks = toks[:1]
		} else {
			rhs += s
			toks = toks[:1]
		}
	}
	return lhs, rhs
}

func checkHs(expr []string, rule string, set string, hs string) error {
	estr := strings.Join(expr, " ")
	if rule == "" {
		err := "empty %s handside in rule `%s`"
		return fmt.Errorf("error: "+err, hs, estr)
	}
	for _, c := range rule {
		if c >= 'A' && c <= 'Z' || inSet(c, set) {
			continue
		} else {
			err := "unknown literal `%c` in %s handside of rule `%s`"
			return fmt.Errorf("error: "+err, c, hs, estr)
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

func toRPN(src []string, expr string, hs string) (string, error) {
	prec := map[rune]int{
		'^': 1,
		'|': 2,
		'+': 4,
		'!': 8,
	}
	estr := strings.Join(src, " ")
	stack := []rune{}
	queue := []rune{}
	for _, c := range expr {
		if c >= 'A' && c <= 'Z' {
			queue = append(queue, c)
		} else if c == '!' || c == '+' || c == '|' || c == '^' {
			for len(stack)-1 > 0 && prec[stack[len(stack)-1]] > prec[c] {
				queue = append(queue, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)
		} else if c == '(' {
			stack = append(stack, c)
		} else if c == ')' {
			for len(stack)-1 > 0 && stack[len(stack)-1] != '(' {
				queue = append(queue, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				err := "error: mismatched parentheses in %s handside `%s` of rule `%s`"
				return string(queue), fmt.Errorf(err, hs, expr, estr)
			}
			if stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
			} else {
				err := "error: mismatched parentheses in %s handside `%s` of rule `%s`"
				return string(queue), fmt.Errorf(err, hs, expr, estr)
			}
		} else {
			panic("error: unreachable code")
		}

	}
	for len(stack) > 0 {
		c := stack[len(stack)-1]
		if c == '(' || c == ')' {
			err := "error: mismatched parentheses in %s handside `%s` of rule `%s`"
			return string(queue), fmt.Errorf(err, hs, expr, estr)
		}
		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return string(queue), nil
}

func simplifyRhs(rhs string) string {
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

func makeRule(lhs string, rhs string) []Rule {
	out := []Rule{}
	rrhs := strrev(rhs)
	for _, c := range rrhs {
		out = append(out, Rule{
			id:   c,
			node: makeNode(lhs),
		})
	}
	return out
}

func makeNode(lhs string) TreeNode {
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
		} else { // TODO: handle not
			tmp := &BinaryGate{
				gType: getType(c),
				left:  nil,
				right: nil,
			}
			// TODO: check against zero size stack in both pops
			t1 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			t2 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			tmp.right = t1
			tmp.left = t2

			tree = tmp
			stack = append(stack, tree)
		}
	}
	// TODO: check against empty stack
	tree = stack[len(stack)-1]
	return tree
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
