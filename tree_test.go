package main

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	f := GetFacts()
	f.Reset()
	f.Set('A', true)
	f.Set('B', true)
	f.Set('C', true)

	// !A ^ (B & C) => D
	tree := &BinaryGate{GateXor,
		&UnaryGate{GateNot,
			&Value{'A'}},
		&BinaryGate{GateAnd,
			&Value{'B'},
			&Value{'C'}}}

	value, err := tree.Evaluate()
	if err != nil {
		panic(err) // TODO: address error
	}
	expected := true
	if value != expected {
		t.Errorf("Expect %v from %v. Got %v", expected, tree, value)
	}
}

func TestValue(t *testing.T) {
	f := GetFacts()
	f.Reset()

	f.Set('A', true)
	f.Set('B', false)

	trees := []TreeNode{
		&Value{'A'},
		&Value{'B'},
	}
	expected := []bool{true, false}

	for ii := range trees {
		val, err := trees[ii].Evaluate()
		if err != nil {
			t.Error(err)
		}
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}

func TestNot(t *testing.T) {
	f := GetFacts()
	f.Reset()

	f.Set('A', true)
	f.Set('B', false)

	trees := []TreeNode{
		&UnaryGate{GateNot,
			&Value{'A'}},
		&UnaryGate{GateNot,
			&Value{'B'}},
	}
	expected := []bool{false, true}

	for ii := range trees {
		val, err := trees[ii].Evaluate()
		if err != nil {
			t.Error(err)
		}
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}

func TestAnd(t *testing.T) {
	f := GetFacts()
	f.Reset()
	f.Set('A', true)
	f.Set('B', false)
	f.Set('C', true)

	trees := []TreeNode{
		&BinaryGate{GateAnd,
			&Value{'A'},
			&Value{'B'}},

		&BinaryGate{GateAnd,
			&Value{'A'},
			&Value{'C'}},

		&BinaryGate{GateAnd,
			&Value{'F'},
			&Value{'G'}},
	}
	expected := []bool{false, true, false}

	for ii := range trees {
		val, err := trees[ii].Evaluate()
		if err != nil {
			t.Error(err)
		}
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}

func TestOr(t *testing.T) {
	f := GetFacts()
	f.Reset()

	f.Set('A', true)
	f.Set('B', true)
	f.Set('C', false)
	f.Set('D', false)

	trees := []TreeNode{
		&BinaryGate{GateOr,
			&Value{'A'},
			&Value{'B'}},

		&BinaryGate{GateOr,
			&Value{'B'},
			&Value{'C'}},

		&BinaryGate{GateOr,
			&Value{'C'},
			&Value{'D'}},
	}
	expected := []bool{true, true, false}

	for ii := range trees {
		val, err := trees[ii].Evaluate()
		if err != nil {
			t.Error(err)
		}
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}

func TestXor(t *testing.T) {
	f := GetFacts()
	f.Reset()

	f.Set('A', true)
	f.Set('B', true)
	f.Set('C', false)
	f.Set('D', false)

	trees := []TreeNode{
		&BinaryGate{GateXor,
			&Value{'A'},
			&Value{'B'}},

		&BinaryGate{GateXor,
			&Value{'B'},
			&Value{'C'}},

		&BinaryGate{GateXor,
			&Value{'C'},
			&Value{'D'}},
	}
	expected := []bool{false, true, false}

	for ii := range trees {
		val, err := trees[ii].Evaluate()
		if err != nil {
			t.Error(err)
		}
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}
