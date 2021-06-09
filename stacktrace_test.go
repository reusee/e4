package e4

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestStacktrace(t *testing.T) {
	TestWrapFunc(t, NewStacktrace())

	trace := NewStacktrace()(io.EOF)
	ok, err := regexp.MatchString(
		`\$ e4.stacktrace_test.go:[0-9]+ .*/e4/ e4.TestStacktrace\n&.*\n&.*\nEOF`,
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

func TestDropFrame(t *testing.T) {
	TestWrapFunc(t, DropFrame(func(frmae Frame) bool {
		return false
	}))

	err := NewStacktrace()(io.EOF)
	err = DropFrame(func(frame Frame) bool {
		return !strings.Contains(frame.Function, "TestDropFrame")
	})(err)
	var stacktrace *Stacktrace
	if !as(err, &stacktrace) {
		t.Fatal()
	}
	if len(stacktrace.Frames) != 1 {
		t.Fatal()
	}
	if !strings.Contains(stacktrace.Frames[0].Function, "TestDropFrame") {
		t.Fatal()
	}
}
