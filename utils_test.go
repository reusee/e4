package e4

import (
	"io"
	"os"
	"testing"
)

type closer func() error

func (c closer) Close() error {
	return c()
}

func TestClose(t *testing.T) {
	err := Wrap(io.EOF,
		Close(closer(func() error {
			return io.ErrClosedPipe
		})),
		Close(closer(func() error {
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

func TestIgnore(t *testing.T) {
	err := func() (err error) {
		defer Handle(&err, Ignore(io.EOF))
		return io.EOF
	}()
	if err != nil {
		t.Fatal()
	}
	err = func() (err error) {
		defer Handle(&err, Ignore(io.EOF))
		return io.ErrClosedPipe
	}()
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
}

func TestIgnoreAs(t *testing.T) {
	err := func() (err error) {
		defer Handle(&err, IgnoreAs(new(*os.PathError)))
		return new(os.PathError)
	}()
	if err != nil {
		t.Fatal()
	}
	err = func() (err error) {
		defer Handle(&err, IgnoreAs(new(*os.PathError)))
		return io.EOF
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
}

func TestIgnoreContains(t *testing.T) {
	err := func() (err error) {
		defer Handle(&err, IgnoreContains("EOF"))
		return io.EOF
	}()
	if err != nil {
		t.Fatal()
	}
	err = func() (err error) {
		defer Handle(&err, IgnoreContains("EOF"))
		return io.ErrClosedPipe
	}()
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
}

func TestUtilFuncs(t *testing.T) {
	testWrapFunc(t, Close(io.NopCloser(os.Stderr)))
	testWrapFunc(t, Do(func() {}))
	testWrapFunc(t, Ignore(io.EOF))
	testWrapFunc(t, IgnoreAs(new(*os.PathError)))
	testWrapFunc(t, IgnoreContains("foo"))
}
