package main

import (
	"fmt"
	"testing"
)

func TestRange(t *testing.T) {
	tests := []byte{'A', 'A' - 1, 'Z', 'Z' + 1, 'G', 'a', 'z'}
	expected := []bool{true, false, true, false, true, false, false}

	f := GetFacts()

	for ii := range tests {
		if v := f.InRange(tests[ii]); v != expected[ii] {
			t.Errorf("Expected %v for value %c.  Got %v", expected[ii], tests[ii], v)
		}
	}
}

func TestSetQuery(t *testing.T) {
	tests := []byte{'A', 'A' - 1, 'Z', 'Z' + 1, 'G', 'a'}
	shouldErr := []bool{false, true, false, true, false, true}
	set := []bool{true, false, false, true, false, false}

	f := GetFacts()
	f.Reset()

	for ii := range tests {
		err := f.Set(tests[ii], set[ii])
		if shouldErr[ii] && err == nil {
			t.Errorf("Should have errored when setting '%c'", tests[ii])
		} else if !shouldErr[ii] && err != nil {
			t.Errorf("Should not have errored when setting '%c': %v",
				tests[ii], err)
		}

		res, err := f.Query(tests[ii])
		if shouldErr[ii] && err == nil {
			t.Errorf("Should have errored when querying '%c'", tests[ii])
		} else if !shouldErr[ii] && err != nil {
			t.Errorf("Should not have errored when querying '%c': %v",
				tests[ii], err)
		}
		if !shouldErr[ii] && res != set[ii] {
			t.Errorf("Did not set '%c' to %v", tests[ii], res)
		}
	}

	issetTests := []byte{'A', 'B', 'F', 'G', 'H'}
	issetBool := []bool{true, false, false, true, false}

	for ii := range issetTests {
		truth, err := f.IsSet(issetTests[ii])
		if err != nil {
			panic(err)
		}
		if truth != issetBool[ii] {
			t.Errorf("Var %c IsSet should be %v.  Was %v", issetTests[ii], issetBool[ii], truth)
		}
	}

	err := f.Set('A', true)
	if err == nil {
		t.Errorf("Should have errored when setting 'A' twice")
	}
}

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
	f.Reset()
	f.UserSet([]byte{'A'})

	for ii := range trees {
		f.AddRule('D', trees[ii])
	}
	f.AddRule('G', &Value{'Z'})
	res, _ := f.UserQuery([]byte{'D'})
	expected := []bool{false}
	for ii := range expected {
		if expected[ii] != res[ii] {
			t.Error("ERR")
		}
	}
	f.AddRule('B', &Value{'A'})
	res, _ = f.UserQuery([]byte{'D', 'A', 'B', 'G'})
	expected = []bool{true, true, true, false}
	for ii := range expected {
		if expected[ii] != res[ii] {
			t.Error("ERR")
		}
	}
}

func TestLongDefinitions(t *testing.T) {
	f := GetFacts()
	f.Reset()
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
