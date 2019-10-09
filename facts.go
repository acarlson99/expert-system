package main

import "fmt"

type Facts struct {
	f [26]bool
}

var g_facts *Facts

func GetFacts() *Facts {
	if g_facts != nil {
		return g_facts
	}
	f := new(Facts)
	f.Init()
	g_facts = f
	return f
}

func (f *Facts) Init() {
	for ii := range f.f {
		f.f[ii] = false
	}
}

func (f *Facts) Reset() {
	for ii := range f.f {
		f.f[ii] = false
	}
}

func (f *Facts) InRange(c byte) bool {
	idx := c - 'A'
	return (idx >= 0 && idx < 26)
}

func (f *Facts) Query(c byte) (bool, error) {
	if !f.InRange(c) {
		return false, fmt.Errorf("Variable '%c' not available", c)
	}
	return f.f[c-'A'], nil
}

func (f *Facts) Set(c byte, t bool) error {
	if !f.InRange(c) {
		return fmt.Errorf("Variable '%c' not available", c)
	}
	f.f[c-'A'] = t
	return nil
}
