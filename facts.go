package main

import "fmt"

type Fact int

const (
	F Fact = 0
	T Fact = 1
	U Fact = 2
)

type Facts struct {
	f [26]Fact
}

func (t Fact) String() string {
	switch t {
	case F:
		return "F"
	case T:
		return "T"
	case U:
		return "U"
	}
	panic(t)
}

func NewFacts() Facts {
	var f Facts
	f.Init()
	return f
}

func (f *Facts) Init() {
	for ii := range f.f {
		f.f[ii] = U
	}
}

func (f *Facts) InRange(c byte) bool {
	idx := c - 'A'
	return (idx >= 0 && idx < 26)
}

func (f *Facts) Query(c byte) (Fact, error) {
	if !f.InRange(c) {
		return U, fmt.Errorf("Variable '%c' not available", c)
	}
	return f.f[c-'A'], nil
}

func (f *Facts) Set(c byte, t Fact) error {
	if !f.InRange(c) {
		return fmt.Errorf("Variable '%c' not available", c)
	}
	f.f[c-'A'] = t
	return nil
}
