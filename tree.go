package main

import "fmt"

type TreeNode interface {
	Evaluate() bool
}

// value literal.  'A', 'B' evaluate to themselves or the child node
type Value struct {
	index byte
	next  TreeNode
}

func (v *Value) Evaluate() bool {
	facts := GetFacts()

	if v.next != nil {
		truth := v.next.Evaluate()
		facts.Set(v.index, truth)
	}
	q, _ := facts.Query(v.index)
	if verbose {
		fmt.Printf("%c evaluating to %v\n", v.index, q)
	}
	return q
}

// Unary gate
type GType int

const (
	GateNot GType = iota
	GateAnd
	GateOr
	GateXor
)

type UnaryGate struct {
	gateType GType
	next     TreeNode
}

func (g *UnaryGate) Evaluate() bool {
	next := g.next
	if next == nil {
		panic("Not operator called on nil ptr")
	}
	nval := next.Evaluate()

	var value bool
	switch g.gateType {
	case GateNot:
		value = !nval
	default:
		panic("Invalid unary gate type: " + string(g.gateType))
	}

	if verbose {
		fmt.Printf("next: %v, evaluating to %v\n", nval, value)
	}
	return value
}

// Binary gate
type BinaryGate struct {
	gateType GType
	left     TreeNode
	right    TreeNode
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

	switch g.gateType {
	case GateAnd:
		value = lval && rval
	case GateOr:
		value = lval || rval
	case GateXor:
		value = lval != rval
	default:
		panic("Invalid binary gate type: " + string(g.gateType))
	}

	if verbose {
		fmt.Printf("left: %v, right: %v, evaluating to %v\n", lval, rval, value)
	}
	return value
}
