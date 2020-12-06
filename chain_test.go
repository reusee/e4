package e4

import (
	"io"
	"testing"
)

func TestChain(t *testing.T) {
	err := Chain{
		Err: Chain{
			Err: io.EOF,
			Prev: Chain{
				Err:  io.ErrClosedPipe,
				Prev: io.ErrNoProgress,
			},
		},
		Prev: Chain{
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
