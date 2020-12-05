package e4

import (
	"io"
	"regexp"
	"testing"
)

func TestStacktrace(t *testing.T) {
	trace := NewStacktrace()(io.EOF)
	ok, err := regexp.MatchString(
		"> at .*stacktrace_test.go:[0-9]+ github.com/reusee/e4.TestStacktrace\n-.*\n-.*\nEOF",
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
