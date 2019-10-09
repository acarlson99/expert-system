package main

import (
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
}
