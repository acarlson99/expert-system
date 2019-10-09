package main

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	f := GetFacts()
	f.Set('A', true)
	f.Set('B', true)
	f.Set('C', true)

	// !A ^ (B & C) => D
	tree := &Value{'D', &BinaryGate{GateXor,
		&UnaryGate{GateNot, &Value{'A', nil}},
		&BinaryGate{GateAnd, &Value{'B', nil}, &Value{'C', nil}}}}

	// verbose = true
	value := tree.Evaluate()
	expected := true
	if value != expected {
		t.Errorf("Expect %v got %v", expected, value)
	}
}
