package catcher

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestBase(t *testing.T) {
	s := Wrap(1, errors.New(""))
	fmt.Println(s)
	fmt.Println(FromError(s))
}

func TestStack(t *testing.T) {
	s := Wrap(1, errors.New("bad"))
	stack, ok := StackFromError(s)
	if !ok {
		t.Fatal("bad...")
	}
	fmt.Println(stack)
	x, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stack.Rel("github.com/Myriad-Dreamin/catcher", x))
}

func TestDescriptor(t *testing.T) {
	SetCodeDescriptor(func(code int) string {
		switch code {
		case CodeDeserializeError:
			return "DeserializeError(-1)"
		default:
			return ErrC(code)
		}
	})

	s := Wrap(1, trace2())
	stack, ok := StackFromError(s)
	if !ok {
		t.Fatal("bad...")
	}
	x, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stack.Rel("github.com/Myriad-Dreamin/catcher", x))
}

func ErrC(i int) string {
	switch i {
	case 1:
		return "Bad"
	case 2:
		return "DecErr"
	default:
		return fmt.Sprintf("Unknown(%v)", i)
	}
}

func trace1() error {
	return Wrap(1, errors.New("hh nb"))
}

func trace2() error {
	return Wrap(2, trace1())
}

func TestStack2(t *testing.T) {
	s := Wrap(1, trace2())
	stack, ok := StackFromError(s)
	if !ok {
		t.Fatal("bad...")
	}
	fmt.Println(stack)
	x, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stack.Rel("github.com/Myriad-Dreamin/catcher", x))
}

func BenchmarkWrap(b *testing.B) {
	err := errors.New("")
	SetErrorFlag(Debug)
	for i := 0; i < b.N; i++ {
		_ = Wrap(1, err)
	}
}

func BenchmarkWrapWithOutCollectInfo(b *testing.B) {
	err := errors.New("")
	SetErrorFlag(Prod)
	for i := 0; i < b.N; i++ {
		_ = Wrap(1, err)
	}
}

func BenchmarkWrapChain(b *testing.B) {
	err := errors.New("")
	SetErrorFlag(Debug)
	for i := 0; i < b.N; i++ {
		_ = Wrap(1, Wrap(1, Wrap(1, err)))
	}
}

func BenchmarkWrapChainWithOutCollectInfo(b *testing.B) {
	err := errors.New("")
	SetErrorFlag(Prod)
	for i := 0; i < b.N; i++ {
		_ = Wrap(1, Wrap(1, Wrap(1, err)))
	}
}

func BenchmarkStackFromError(b *testing.B) {
	stack := Wrap(1, Wrap(1, Wrap(1, errors.New("")))).Error()
	SetErrorFlag(Debug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = StackFromString(stack)
	}
}

func BenchmarkStackFromErrorWithOutCollectInfo(b *testing.B) {
	stack := Wrap(1, Wrap(1, Wrap(1, errors.New("")))).Error()
	SetErrorFlag(Prod)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = StackFromString(stack)
	}
}
