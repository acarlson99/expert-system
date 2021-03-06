package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type Fact struct {
	truth       bool
	visited     bool
	userdefined bool
	rule        TreeNode
}

func (fact *Fact) Evaluate() bool {
	if fact.rule == nil {
		return fact.truth
	}
	return fact.rule.Evaluate()
}

func (fact *Fact) Query() bool {
	if fact.truth || fact.visited || fact.rule == nil {
		return fact.truth
	}
	fact.visited = true
	val := fact.Evaluate()
	fact.truth = val
	return fact.truth
}

type Facts struct {
	f [26]Fact
}

var g_facts *Facts

func GetFacts() *Facts {
	if g_facts != nil {
		return g_facts
	}
	f := new(Facts)
	f.HardReset()
	g_facts = f
	return f
}

func (f *Facts) HardReset() {
	for ii := range f.f {
		f.f[ii] = Fact{false, false, false, nil}
	}
}

func (f *Facts) SoftReset() {
	for ii := range f.f {
		if !f.f[ii].userdefined {
			f.f[ii] = Fact{false, false, false, f.f[ii].rule}
		} else {
			f.f[ii].visited = false
		}
	}
}

func (f *Facts) UserSet(cs []byte) {
	for ii := range f.f {
		f.f[ii] = Fact{false, false, false, f.f[ii].rule}
	}
	for _, c := range cs {
		f.f[c-'A'] = Fact{true, false, true, f.f[c-'A'].rule}
	}
}

func (f *Facts) UserQuery(cs []byte) []bool {
	f.SoftReset()
	res := []bool{}
	for _, c := range cs {
		if verbose {
			fmt.Printf("Querying %c\n", c)
		}
		val := f.f[c-'A'].Query()
		res = append(res, val)
	}
	return res
}

func (f *Facts) AddRule(c byte, t TreeNode) {
	f.SoftReset()
	fact := &f.f[c-'A']
	if fact.rule == nil {
		fact.rule = t
	} else {
		constructed := &BinaryGate{GateOr, t, fact.rule}
		fact.rule = constructed
	}
}

func (f *Facts) Get(c byte) *Fact {
	return &f.f[c-'A']
}

func (f *Facts) ToGraphviz() *gographviz.Graph {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	f.SoftReset()
	for ii, _ := range f.f {
		fact := &f.f[ii]
		if fact.rule != nil && !fact.visited {
			name := string(ii + 'A')
			fact.visited = true
			value, nname := fact.rule.AddToGraph(graph)
			m := make(map[string]string)
			if fact.userdefined {
				m["color"] = "blue"
			} else if value {
				m["color"] = "green"
			}
			if !fact.userdefined {
				fact.truth = value
			}
			graph.AddNode("G", name, m)
			graph.AddEdge(name, nname, true, nil)
			fact.visited = true
		}
	}
	return graph
}
