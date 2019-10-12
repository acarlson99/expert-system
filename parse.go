package main

import (
	"fmt"
	"strings"
)

func Parse(src string) (interface{}, error) {
	if len(src) == 0 {
		return nil, nil
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
	} else {
		err := "error: unknown expression `%s`"
		return nil, fmt.Errorf("error: "+err, src)
	}
	return nil, nil
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

type Query []byte

func query(src string) (Query, error) {
	out := []byte{}
	if len(src) == 1 {
		return out, nil
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

func parseRule(src string, lhs string, rhs string) ([]Rule, error) {
	out := []Rule{}
	if err := checkHs(src, lhs, "()!+|^", "left"); err != nil {
		return out, err
	}
	if err := checkHs(src, rhs, "()+", "right"); err != nil {
		return out, err
	}
	lhs1, err1 := toRPN(src, lhs, "left")
	if err1 != nil {
		return out, err1
	}
	rhs1, err2 := toRPN(src, rhs, "left")
	if err2 != nil {
		return out, err2
	}
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

func toRPN(src string, rule string, hs string) (string, error) {
	prec := map[rune]int{
		'^': 1,
		'|': 2,
		'+': 4,
		'!': 8,
	}
	stack := []rune{}
	queue := []rune{}
	for _, c := range rule {
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
				return string(queue), fmt.Errorf(err, hs, rule, src)
			}
			if stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
			} else {
				err := "error: mismatched parentheses in %s handside `%s` of rule `%s`"
				return string(queue), fmt.Errorf(err, hs, rule, src)
			}
		} else {
			err := "unknown literal `%c` in `%s`"
			return string(queue), fmt.Errorf("i-error: "+err, c, "toRPN")
		}

	}
	for len(stack) > 0 {
		c := stack[len(stack)-1]
		if c == '(' || c == ')' {
			err := "error: mismatched parentheses in %s handside `%s` of rule `%s`"
			return string(queue), fmt.Errorf(err, hs, rule, src)
		}
		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return string(queue), nil
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
			t1 = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			tmp.next = t1
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
