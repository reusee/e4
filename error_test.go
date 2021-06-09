package e4

import (
	"io"
	"testing"
)

func TestError(t *testing.T) {
	err := MakeErr(
		MakeErr(
			io.EOF,
			MakeErr(
				io.ErrClosedPipe,
				io.ErrNoProgress,
			),
		),
		MakeErr(
			io.ErrShortBuffer,
			io.ErrUnexpectedEOF,
		),
	)
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

func TestWith(t *testing.T) {
	testWrapFunc(t, With(io.EOF))
}
