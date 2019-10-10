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
	tree := &Value{'D',
		&BinaryGate{GateXor,
			&UnaryGate{GateNot,
				&Value{'A', nil}},
			&BinaryGate{GateAnd,
				&Value{'B', nil},
				&Value{'C', nil}}}}

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
		&Value{'A', nil},
		&Value{'B', nil},
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
			&Value{'A', nil}},
		&UnaryGate{GateNot,
			&Value{'B', nil}},
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
		&Value{'D',
			&BinaryGate{GateAnd,
				&Value{'A', nil},
				&Value{'B', nil}}},

		&Value{'E',
			&BinaryGate{GateAnd,
				&Value{'A', nil},
				&Value{'C', nil}}},

		&Value{'H',
			&BinaryGate{GateAnd,
				&Value{'F', nil},
				&Value{'G', nil}}},
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
		&Value{'E',
			&BinaryGate{GateOr,
				&Value{'A', nil},
				&Value{'B', nil}}},

		&Value{'F',
			&BinaryGate{GateOr,
				&Value{'B', nil},
				&Value{'C', nil}}},

		&Value{'G',
			&BinaryGate{GateOr,
				&Value{'C', nil},
				&Value{'D', nil}}},
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
		&Value{'E',
			&BinaryGate{GateXor,
				&Value{'A', nil},
				&Value{'B', nil}}},

		&Value{'F',
			&BinaryGate{GateXor,
				&Value{'B', nil},
				&Value{'C', nil}}},

		&Value{'G',
			&BinaryGate{GateXor,
				&Value{'C', nil},
				&Value{'D', nil}}},
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
