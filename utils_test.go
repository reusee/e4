package e4

import (
	"io"
	"testing"
)

type closer func() error

func (c closer) Close() error {
	return c()
}

func TestWithCloser(t *testing.T) {
	err := Wrap(io.EOF,
		WithClose(closer(func() error {
			return io.ErrClosedPipe
		})),
		WithClose(closer(func() error {
			return nil
		})),
	)
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
}
