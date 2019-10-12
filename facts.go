package main

import "fmt"

type Fact struct {
	// value of Fact.  true/false
	t bool
	// whether or not Fact has been set or implicit
	set bool
	// rule
	rule TreeNode
	// user defined or inferred by system
	userdefined bool
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
		f.f[ii] = Fact{false, false, nil, false}
	}
}

func (f *Facts) SoftReset() {
	for ii := range f.f {
		if !f.f[ii].userdefined {
			f.f[ii] = Fact{false, false, f.f[ii].rule, false}
		}
	}
}

func (f *Facts) UserSet(cs []byte) error {
	f.SoftReset()
	for ii := range cs {
		if !f.InRange(cs[ii]) {
			return fmt.Errorf("Variable '%c' not available", cs[ii])
		}
		f.f[cs[ii]-'A'] = Fact{true, true, nil, true}
	}
	return nil
}

func (f *Facts) InRange(c byte) bool {
	idx := c - 'A'
	return (idx >= 0 && idx < 26)
}

func (f *Facts) Query(c byte) (bool, error) {
	if !f.InRange(c) {
		return false, fmt.Errorf("Variable '%c' not available", c)
	}
	fact := &f.f[c-'A']
	if a, _ := f.IsSet(c); a {
		return fact.t, nil
	} else if fact.rule != nil {
		err := f.Evaluate(c)
		if err != nil {
			return false, err
		}
		return fact.t, nil
	} else {
		return fact.t, nil
	}
}

func (f *Facts) UserQuery(cs []byte) ([]bool, error) {
	f.SoftReset()
	res := []bool{}
	for ii := range cs {
		if !f.InRange(cs[ii]) {
			return res, fmt.Errorf("Variable '%c' not available", cs[ii])
		}
		if verbose { // TODO: kill - don't be so mean to your code :c
			fmt.Println("EVALUATING")
		}
		if err := f.Evaluate(cs[ii]); err != nil {
			return res, err
		}
		q, _ := f.Query(cs[ii])
		res = append(res, q)
	}
	return res, nil
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

func (f *Facts) AddRule(c byte, t TreeNode) error {
	if !f.InRange(c) {
		return fmt.Errorf("Variable '%c' not available", c)
	} else if t == nil {
		return fmt.Errorf("Assigning nil rule to variable '%c'", c)
	}
	fact := &f.f[c-'A']
	if fact.rule == nil {
		fact.rule = t
	} else {
		constructed := &BinaryGate{GateOr, t, fact.rule}
		fact.rule = constructed
	}
	return nil
}

func (f *Facts) Evaluate(c byte) error {
	if !f.InRange(c) {
		return fmt.Errorf("Variable '%c' not available", c)
	}
	fact := &f.f[c-'A']
	if fact.rule == nil {
		return nil
	}
	if verbose { // TODO: kill
		fmt.Println("CALLING EVALUATE ON", c)
	}
	value := fact.rule.Evaluate()
	if !fact.set && value {
		return f.Set(c, value)
	}
	if fact.set && fact.t != value {
		return fmt.Errorf("Variable '%c' set to different values", c)
	}
	return nil
}
