package e4

import (
	"io"
	"os"
	"testing"
)

func TestWrapFunc(t *testing.T) {

	fn := WrapFunc(func(err error) error {
		return Chain{
			Err: io.EOF,
			Prev: Chain{
				Err: io.ErrClosedPipe,
				Prev: Chain{
					Err:  new(os.PathError),
					Prev: err,
				},
			},
		}
	})
	err, ok := any(fn).(error)
	if !ok {
		t.Fatal()
	}
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
	var pathError *os.PathError
	if !as(err, &pathError) {
		t.Fatal()
	}

}
