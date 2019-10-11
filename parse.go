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

func parseUnop(op string) ([]rune, error) {
	switch op[0] {
	case '=':
		return assign(op[1:])
	case '?':
		return query(op[1:])
	default:
		// TODO: handle builtins
		return nil, nil
	}
}

func assign(src string) ([]rune, error) {
	out := []rune{}
	if len(src) == 0 {
		err := "empty left handside in assignment"
		return out, fmt.Errorf("error: " + err)
	}
	for _, c := range src {
		fmt.Println(c)
		if c >= 'A' && c <= 'Z' {
			out = append(out, c)
		} else {
			err := "unknown literal `%c` in assignment `%s`"
			return out, fmt.Errorf("error: "+err, c, src)
		}
	}
	return out, nil
}

func query(src string) ([]rune, error) {
	out := []rune{}
	if len(src) == 0 {
		err := "empty left handside in query"
		return out, fmt.Errorf("error: " + err)
	}
	for _, c := range src {
		fmt.Println(c)
		if c >= 'A' && c <= 'Z' {
			out = append(out, c)
		} else {
			err := "unknown literal `%c` in query `%s`"
			return out, fmt.Errorf("error: "+err, c, src)
		}
	}
	return out, nil
}

// @rule

func parseRule(toks []string) (string, error) {
	lhs, rhs := splitHs(toks)
	if err := checkHs(toks, lhs, "()!+|^", "left"); err != nil {
		return "", err
	}
	if err := checkHs(toks, rhs, "()+", "right"); err != nil {
		return "", err
	}
	lhs1, err1 := toRPN(toks, lhs, "left")
	if err1 != nil {
		return "", err1
	}
	rhs1, err2 := toRPN(toks, rhs, "left")
	if err2 != nil {
		return "", err2
	}
	return lhs1 + rhs1, nil
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
