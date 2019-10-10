package main

import "fmt"

type TreeNode interface {
	Evaluate() bool
	String() string
}

// value literal.  'A', 'B' evaluate to their boolean values
type Value struct {
	ch byte
}

func (v *Value) Evaluate() bool {
	facts := GetFacts()

	value, _ := facts.Query(v.ch)
	if verbose {
		fmt.Printf("%c = %v\n", v.ch, value)
	}
	return value
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

func (g *UnaryGate) Evaluate() bool {
	next := g.next
	if next == nil {
		panic("Not operator called on nil ptr")
	}
	nval := next.Evaluate()

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
	return value
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

func (g *BinaryGate) Evaluate() bool {
	left := g.left
	right := g.right
	if left == nil || right == nil {
		panic("Binary operator called on nil ptr")
	}

	var value bool
	switch g.gType {
	case GateAnd:
		value = left.Evaluate() && right.Evaluate()
	case GateOr:
		value = left.Evaluate() || right.Evaluate()
	case GateXor:
		value = left.Evaluate() != right.Evaluate()
	default:
		panic("Invalid binary gate type: " + string(g.gType))
	}

	if verbose {
		fmt.Printf("%v %s %v = %v\n", left, g.gType, right, value)
	}
	return value
}

func (g *BinaryGate) String() string {
	return "(" + g.left.String() + " " + g.gType.String() + " " + g.right.String() + ")"
}
