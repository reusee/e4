package e4

import (
	"io"
	"testing"
)

func TestWrapAssign(t *testing.T) {
	err := func() (err error) {
		we := Wrap.With(With(io.ErrClosedPipe))
		defer we.Assign(&err)
		return io.EOF
	}()
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
}
