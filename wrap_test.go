package e4

import (
	"io"
	"testing"
)

func TestWrap(t *testing.T) {
	err := NewInfo("foo")(io.EOF)
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if info := new(Info); !as(err, &info) {
		t.Fatal()
	} else if *info != "foo" {
		t.Fatal()
	}
	err = NewInfo("foo")
	if err.Error() != "foo" {
		t.Fatal()
	}
}
