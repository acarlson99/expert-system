package main

import "fmt"

type TreeNode interface {
	Evaluate() (bool, error)
	String() string
}

// value literal.  'A', 'B' evaluate to themselves or the child node
type Value struct {
	ch   byte
	next TreeNode
}

func (v *Value) Evaluate() (bool, error) {
	facts := GetFacts()

	if v.next != nil {
		truth, err := v.next.Evaluate()
		if err != nil {
			return false, err // TODO: address error
		}
		if set, _ := facts.IsSet(v.ch); !set {
			facts.Set(v.ch, truth)
		} else {
			factVal, _ := facts.Query(v.ch)
			if truth != factVal {
				return false, fmt.Errorf("Conflicting definitions for %c", v.ch)
			}
		}
	}
	value, _ := facts.Query(v.ch)
	if verbose {
		fmt.Printf("%v => %c\n", v.next, v.ch)
		fmt.Printf("%c = %v\n", v.ch, value)
	}
	return value, nil
}

func (v *Value) String() string {
	if v.next != nil {
		return v.next.String() + " => " + string(v.ch)
	}
	return string(v.ch)
}

// Gate type enum
type GType int

const (
	GateNot GType = iota
	GateAnd
	GateOr
	GateXor
)

func (t GType) String() string {
	switch t {
	case GateNot:
		return "!"
	case GateAnd:
		return "&"
	case GateOr:
		return "|"
	case GateXor:
		return "^"
	default:
		return "?"
	}
}

// Unary gate

type UnaryGate struct {
	gType GType
	next  TreeNode
}

func (g *UnaryGate) Evaluate() (bool, error) {
	next := g.next
	if next == nil {
		panic("Not operator called on nil ptr")
	}
	nval, err := next.Evaluate()
	if err != nil {
		return false, err // TODO: address error
	}

	var value bool
	switch g.gType {
	case GateNot:
		value = !nval
	default:
		panic("Invalid unary gate type: " + string(g.gType))
	}

	if verbose {
		fmt.Printf("%s%v = %v\n", g.gType, g.next, value)
	}
	return value, nil
}

func (g *UnaryGate) String() string {
	return g.gType.String() + g.next.String()
}

// Binary gate
type BinaryGate struct {
	gType GType
	left  TreeNode
	right TreeNode
}

// TODO: cache results intelligently
func (g *BinaryGate) Evaluate() (bool, error) {
	left := g.left
	right := g.right
	if left == nil || right == nil {
		panic("Binary operator called on nil ptr")
	}
	lval, err := left.Evaluate()
	if err != nil {
		return false, err // TODO: address error
	}
	rval, err := right.Evaluate()
	if err != nil {
		return false, err
	}
	var value bool

	switch g.gType {
	case GateAnd:
		value = lval && rval
	case GateOr:
		value = lval || rval
	case GateXor:
		value = lval != rval
	default:
		panic("Invalid binary gate type: " + string(g.gType))
	}

	if verbose {
		fmt.Printf("%v %s %v = %v\n", left, g.gType, right, value)
	}
	return value, nil
}

func (g *BinaryGate) String() string {
	return "(" + g.left.String() + " " + g.gType.String() + " " + g.right.String() + ")"
}
