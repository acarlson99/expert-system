package main

import (
	"fmt"
	"testing"
)

func TestAddRule(t *testing.T) {
	trees := []TreeNode{
		// A + B => D
		// F + G => D
		&BinaryGate{GateAnd,
			&Value{'A'},
			&Value{'B'}},
		&BinaryGate{GateAnd,
			&Value{'F'},
			&Value{'G'}},
	}

	f := GetFacts()
	f.HardReset()
	// =A
	f.UserSet([]byte{'A'})

	for ii := range trees {
		f.AddRule('D', trees[ii])
	}
	// Z => G
	f.AddRule('G', &Value{'Z'})
	res := f.UserQuery([]byte{'D'})
	expected := []bool{false}
	for ii := range expected {
		if expected[ii] != res[ii] {
			t.Error("ERR")
		}
	}
	// A => B
	f.AddRule('B', &Value{'A'})
	res = f.UserQuery([]byte{'D', 'A', 'B', 'G'})
	expected = []bool{true, true, true, false}
	for ii := range expected {
		if expected[ii] != res[ii] {
			t.Error("ERR")
		}
	}
}

func TestLongDefinitions(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	f.AddRule('A', &Value{'B'})
	f.AddRule('B', &Value{'C'})
	f.AddRule('C', &Value{'D'})
	verbose = true
	fmt.Println("START")
	fmt.Println(f.UserQuery([]byte{'A'}))
	fmt.Println("END")
	f.UserSet([]byte{'D'})
	fmt.Println(f.UserQuery([]byte{'A'}))
	fmt.Println("A")
	fmt.Println(f.UserQuery([]byte{'D'}))
	verbose = false
}
