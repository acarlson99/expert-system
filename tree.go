package main

import "fmt"

type TreeNode interface {
	Evaluate() (bool, error)
	String() string
}

// value literal.  'A', 'B' evaluate to their boolean values
type Value struct {
	ch byte
}

func (v *Value) Evaluate() (bool, error) {
	facts := GetFacts()

	value, _ := facts.Query(v.ch)
	if verbose {
		fmt.Printf("%c = %v\n", v.ch, value)
	}
	return value, nil
}

func (v *Value) String() string {
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
		return false, err
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

func (g *BinaryGate) Evaluate() (bool, error) {
	left := g.left
	right := g.right
	if left == nil || right == nil {
		panic("Binary operator called on nil ptr")
	}
	lval, err := left.Evaluate()
	if err != nil {
		return false, err
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
