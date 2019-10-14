package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type TreeNode interface {
	Evaluate() bool
	String() string
	// return name of node created
	AddToGraph(graph *gographviz.Graph) (bool, string)
	Name() string
}

// value literal.  'A', 'B' evaluate to their boolean values
type Value struct {
	ch byte
}

func (v *Value) Evaluate() bool {
	facts := GetFacts()

	value := facts.Get(v.ch).Query()
	if verbose {
		fmt.Printf("%c = %v\n", v.ch, value)
	}
	return value
}

func (v *Value) String() string {
	return string(v.ch)
}

func (v *Value) AddToGraph(graph *gographviz.Graph) (bool, string) {
	m := make(map[string]string)
	name := v.Name()
	facts := GetFacts()
	fact := facts.Get(v.ch)

	var value bool
	var nname string
	if fact.visited {
		return fact.truth, name
	}
	if fact.rule != nil && !fact.visited {
		value, nname = fact.rule.AddToGraph(graph)
	} else {
		value = fact.truth
		nname = ""
	}
	if !fact.truth {
		fact.truth = value
	}
	fact.visited = true

	if fact.userdefined {
		m["color"] = "blue"
	} else if fact.truth {
		m["color"] = "green"
	}
	graph.AddNode("G", name, m)
	if nname != "" {
		graph.AddEdge(name, nname, true, nil)
	}
	return fact.truth, name
}

func (v *Value) Name() string {
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
		return "+"
	case GateOr:
		return "|"
	case GateXor:
		return "^"
	default:
		return "?"
	}
}

func (t GType) Word() string {
	switch t {
	case GateNot:
		return "Not"
	case GateAnd:
		return "And"
	case GateOr:
		return "Or"
	case GateXor:
		return "Xor"
	default:
		return "Unknown"
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

func (g *UnaryGate) AddToGraph(graph *gographviz.Graph) (bool, string) {
	next := g.next
	if next == nil {
		panic("Not operator called on nil ptr")
	}
	nval, nname := next.AddToGraph(graph)

	var value bool
	switch g.gType {
	case GateNot:
		value = !nval
	default:
		panic("Invalid unary gate type: " + string(g.gType))
	}

	name := g.Name()
	m := make(map[string]string)
	if value {
		m["color"] = "green"
	}
	m["label"] = g.gType.Word()
	graph.AddNode("G", name, m)
	graph.AddEdge(name, nname, true, nil)

	if verbose {
		fmt.Printf("%s%v = %v\n", g.gType, g.next, value)
	}
	return value, name
}

func (g *UnaryGate) Name() string {
	return g.gType.Word() + fmt.Sprintf("_%p", g)
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

func (g *BinaryGate) AddToGraph(graph *gographviz.Graph) (bool, string) {
	left := g.left
	right := g.right
	if left == nil || right == nil {
		panic("Binary operator called on nil ptr")
	}

	lval, lname := left.AddToGraph(graph)
	rval, rname := right.AddToGraph(graph)

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

	m := make(map[string]string)
	if value {
		m["color"] = "green"
	}
	m["label"] = g.gType.Word()

	name := g.Name()
	graph.AddNode("G", name, m)
	graph.AddEdge(name, lname, true, nil)
	graph.AddEdge(name, rname, true, nil)

	if verbose {
		fmt.Printf("%v %s %v = %v\n", left, g.gType, right, value)
	}
	return value, name
}

func (g *BinaryGate) Name() string {
	return g.gType.Word() + fmt.Sprintf("_%p", g)
}
