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

	value := tree.Evaluate()
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
		val := trees[ii].Evaluate()
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
		// !A
		&UnaryGate{GateNot,
			&Value{'A'}},
		// !B
		&UnaryGate{GateNot,
			&Value{'B'}},
	}
	expected := []bool{false, true}

	for ii := range trees {
		val := trees[ii].Evaluate()
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
		// A + B
		&BinaryGate{GateAnd,
			&Value{'A'},
			&Value{'B'}},

		// A + C
		&BinaryGate{GateAnd,
			&Value{'A'},
			&Value{'C'}},

		// F + G
		&BinaryGate{GateAnd,
			&Value{'F'},
			&Value{'G'}},
	}
	expected := []bool{false, true, false}

	for ii := range trees {
		val := trees[ii].Evaluate()
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
		// A | B
		&BinaryGate{GateOr,
			&Value{'A'},
			&Value{'B'}},

		// B | C
		&BinaryGate{GateOr,
			&Value{'B'},
			&Value{'C'}},

		// C | D
		&BinaryGate{GateOr,
			&Value{'C'},
			&Value{'D'}},
	}
	expected := []bool{true, true, false}

	for ii := range trees {
		val := trees[ii].Evaluate()
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
		// A ^ B
		&BinaryGate{GateXor,
			&Value{'A'},
			&Value{'B'}},

		// B ^ C
		&BinaryGate{GateXor,
			&Value{'B'},
			&Value{'C'}},

		// C ^ D
		&BinaryGate{GateXor,
			&Value{'C'},
			&Value{'D'}},
	}
	expected := []bool{false, true, false}

	for ii := range trees {
		val := trees[ii].Evaluate()
		if val != expected[ii] {
			t.Errorf("Expect %v from %v. Got %v", expected[ii], trees[ii], val)
		}
	}
}
