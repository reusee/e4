package e4

import (
	"errors"
	"io"
	"regexp"
	"testing"
)

func TestStacktrace(t *testing.T) {
	trace := NewStacktrace()(io.EOF)
	ok, err := regexp.MatchString(
		`\$ stacktrace_test.go:[0-9]+ .*/e4/ e4.TestStacktrace\n&.*\n&.*\nEOF`,
		trace.Error(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("got %s", trace.Error())
	}
	if !is(trace, io.EOF) {
		t.Fatal()
	}
}

func TestDeepStacktrace(t *testing.T) {
	var foo func(int) error
	foo = func(i int) error {
		if i < 128 {
			return foo(i + 1)
		}
		return NewStacktrace()(io.EOF)
	}
	err := foo(1)
	if !errors.Is(err, io.EOF) {
		t.Fatal()
	}
}
