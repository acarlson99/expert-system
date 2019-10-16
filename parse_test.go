package main

import (
	"testing"
)

// Recursive Definitions
func Test0(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A=>A")
	if err == nil {
		t.Errorf("recursive definition")
	}
}

func Test1(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A+B|C|D|E|F+G+(Z+Z)+!!!!!!!!!!Z=>Z")
	if err == nil {
		t.Errorf("recursive definition")
	}
}

// Invalid expression
func Test2(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("=>A")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test3(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A=>")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test4(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A+=>B")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test5(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A-=>B")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test6(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("ABC=>E")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test7(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A|B|C=>ED")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test8(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A|B|C=>E!D")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

func Test9(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("A|B|C=>!")
	if err == nil {
		t.Errorf("invalid expression")
	}
}

// success
func TestA(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("(((A)))=>(((Z)))")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestB(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("(((!!!A)|A|A|A+A))=>(((Z)+Z+Z+(Z)))")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestC(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("Z+Z+Z+Z+Z+Z+!!!!!!!!!Z=>((A)+((B)))+C")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

// failure
func TestD(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("()=>()")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func TestE(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("((()))=>((()))")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func TestF(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("((()))()()=>((()))()()")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func Test10(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("(((!!!!Z)))=>(((+)))")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}
