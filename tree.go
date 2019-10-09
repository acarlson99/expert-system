package main

import "fmt"

type TreeNode interface {
	Evaluate() bool
	String() string
}

// value literal.  'A', 'B' evaluate to themselves or the child node
type Value struct {
	ch   byte
	next TreeNode
}

func (v *Value) Evaluate() bool {
	facts := GetFacts()

	if v.next != nil {
		truth := v.next.Evaluate()
		facts.Set(v.ch, truth)
	}
	q, _ := facts.Query(v.ch)
	if verbose {
		fmt.Printf("%c evaluating to %v\n", v.ch, q)
	}
	return q
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
		fmt.Printf("next: %v, evaluating to %v\n", nval, value)
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

// TODO: cache results intelligently
func (g *BinaryGate) Evaluate() bool {
	left := g.left
	right := g.right
	if left == nil || right == nil {
		panic("Binary operator on nil ptr")
	}
	lval := left.Evaluate()
	rval := right.Evaluate()
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
		fmt.Printf("left: %v, right: %v, evaluating to %v\n", lval, rval, value)
	}
	return value
}

func (g *BinaryGate) String() string {
	return "(" + g.left.String() + ") " + g.gType.String() + " (" + g.right.String() + ")"
}
