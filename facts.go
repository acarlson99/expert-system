package main

import "fmt"

type Fact struct {
	// value of Fact.  true/false
	t bool
	// whether or not Fact has been set or implicit
	set bool
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
	f.Reset()
	g_facts = f
	return f
}

func (f *Facts) Reset() {
	for ii := range f.f {
		f.f[ii] = Fact{false, false}
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
	return f.f[c-'A'].t, nil
}

func (f *Facts) IsSet(c byte) (bool, error) {
	if !f.InRange(c) {
		return false, fmt.Errorf("Variable '%c' not available", c)
	}
	return f.f[c-'A'].set, nil
}

func (f *Facts) Set(c byte, t bool) error {
	if !f.InRange(c) {
		return fmt.Errorf("Variable '%c' not available", c)
	}
	fact := &f.f[c-'A']
	if fact.set {
		return fmt.Errorf("Variable '%c' already set", c)
	}
	fact.t = t
	fact.set = true
	return nil
}
