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

// Oneshot
// Success
func Test11(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("?=A|B|C")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func Test12(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("?=(A|(B|(C)))")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func Test13(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("?=(A|(B|(C)))+A+A+A+A")
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

// failure
func Test14(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("?=(A|(B|(C)))+A+A+A+A=>D")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func Test15(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	_, err := Parse("?=(A|(B|(C)))+A+A+A+Aa")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func Test16(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=((A|(B|(C)))+A+A+A+A)!A!")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test17(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=A!")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test18(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=A+B!A")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test19(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=A+B!+A")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

//success
func Test1A(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=A+B+!A")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

func Test1B(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!!!!!Z+!(A)")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

func Test1C(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!!!!!Z+!(A|B|C)")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

func Test1D(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!!!!!Z+!(Z+(A|B|C))")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

func Test1E(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!A+!!!!A")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

func Test1F(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=(A)+(((A)))+!!!(A+(Z|(X)))^Z")
	if err != nil {
		t.Errorf("err: %v", ret)
	}
}

// failure
func Test20(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!()")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test21(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!()A")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test22(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!()A!")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test23(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=!(())")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test24(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("?=(())!")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test25(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("+AB=>C")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test26(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("AB+=>C")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test27(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("AB++=>C")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test28(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("A+B=>C+")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}

func Test29(t *testing.T) {
	f := GetFacts()
	f.HardReset()
	ret, err := Parse("A+B=>+C")
	if err == nil {
		t.Errorf("err: %v", ret)
	}
}
