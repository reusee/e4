package e4

import (
	"io"
	"testing"
)

func TestError(t *testing.T) {
	err := Error{
		Err: Error{
			Err: io.EOF,
			Prev: Error{
				Err:  io.ErrClosedPipe,
				Prev: io.ErrNoProgress,
			},
		},
		Prev: Error{
			Err:  io.ErrShortBuffer,
			Prev: io.ErrUnexpectedEOF,
		},
	}
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}
	if !is(err, io.ErrShortBuffer) {
		t.Fatal()
	}
	if !is(err, io.ErrUnexpectedEOF) {
		t.Fatal()
	}
}
