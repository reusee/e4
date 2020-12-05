package e4

import (
	"io"
	"testing"
)

func TestStack(t *testing.T) {
	err := func() (err error) {
		defer Handle(&err)
		Check(io.EOF, Stack(func() error {
			return io.ErrShortWrite
		}))
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrShortWrite) {
		t.Fatal()
	}
}
