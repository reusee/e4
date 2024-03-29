package e4

import (
	"errors"
	"io"
	"regexp"
	"testing"
)

var (
	is = errors.Is
	as = errors.As
)

func TestCheck(t *testing.T) {

	// check and handle
	err := func() (err error) {
		defer Handle(&err)
		Check(io.EOF)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}

	// check with wrap funcs
	err = func() (err error) {
		defer Handle(&err)
		Check(io.EOF, Info("foo %s", "bar"))
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	ok, e := regexp.MatchString(
		"foo bar\nEOF",
		err.Error(),
	)
	if e != nil {
		t.Fatal(e)
	}
	if !ok {
		t.Fatalf("got %s", err.Error())
	}

	// handle with wrap funcs
	err = func() (err error) {
		defer Handle(&err, Info("foo %s", "bar"))
		Check(io.EOF)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	ok, e = regexp.MatchString(
		`foo bar\nEOF`,
		err.Error(),
	)
	if e != nil {
		t.Fatal(e)
	}
	if !ok {
		t.Fatalf("got %s", err.Error())
	}

	// check and handle nil
	func() {
		defer func() {
			p := recover()
			te, ok := p.(*throw)
			if !ok {
				t.Fatal()
			}
			if !is(te.err, io.EOF) {
				t.Fatal()
			}
		}()
		func() (err error) {
			defer Handle(nil)
			Check(io.EOF)
			return
		}()
	}()

	// non-check panic
	func() {
		defer func() {
			p := recover()
			if p != 42 {
				t.Fatal()
			}
		}()
		func() (err error) {
			defer Handle(&err)
			panic(42)
		}()
	}()

	// return and handle
	err = func() (err error) {
		defer Handle(&err, Info("foo %s", "bar"))
		return io.EOF
	}()
	if err.Error() != "foo bar\nEOF" {
		t.Fatal()
	}
	if !is(err, io.EOF) {
		t.Fatal()
	}

	// return nil and handle
	err = func() (err error) {
		defer Handle(&err, Info("foo %s", "bar"))
		return nil
	}()
	if err != nil {
		t.Fatal()
	}

	// set error to nil in check wrap func
	err = func() (err error) {
		defer Handle(&err)
		Check(io.EOF, WrapFunc(func(err error) error {
			return nil
		}))
		return
	}()
	if err != nil {
		t.Fatal()
	}

	// set error to nil in handle wrap func
	err = func() (err error) {
		defer Handle(
			&err,
			WrapFunc(func(err error) error {
				return nil
			}),
		)
		return io.EOF
	}()
	if err != nil {
		t.Fatal()
	}

	// set return value and check
	err = func() (err error) {
		defer Handle(&err)
		err = io.EOF
		Check(io.ErrNoProgress)
		return
	}()
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}
	if !is(err, io.EOF) {
		t.Fatal()
	}

	// set return value and check the same error
	err = func() (err error) {
		defer Handle(&err)
		err = io.EOF
		Check(err)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}

	// set return value then filter to nil and check
	err = func() (err error) {
		defer Handle(&err)
		err = io.EOF
		Check(io.ErrNoProgress, WrapFunc(func(prev error) error {
			return nil
		}))
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if is(err, io.ErrNoProgress) {
		t.Fatal()
	}

	// auto chain wrap func returns
	err = func() (err error) {
		defer Handle(&err)
		Check(
			io.EOF,
			WrapFunc(func(err error) error {
				return io.ErrClosedPipe
			}),
		)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}

	// multiple handle
	err = func() (err error) {
		defer Handle(&err)
		defer Handle(&err)
		return io.EOF
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}

	// multiple handle and wrap
	err = func() (err error) {
		defer Handle(&err, With(io.ErrClosedPipe))
		defer Handle(&err, With(io.ErrNoProgress))
		return io.EOF
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}

	// check and multiple handle and wrap
	err = func() (err error) {
		defer Handle(&err, With(io.ErrClosedPipe))
		defer Handle(&err, With(io.ErrNoProgress))
		Check(io.EOF)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}

	// Check in defer function
	err = func() (err error) {
		defer Handle(&err)
		defer func() {
			Check(io.EOF)
		}()
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}

}

func TestMust(t *testing.T) {
	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Fatal()
			}
			err, ok := p.(error)
			if !ok {
				t.Fatal()
			}
			if !is(err, io.EOF) {
				t.Fatal()
			}
		}()
		Must(io.EOF)
	}()
}

func TestCheckIgnoreNoHandle(t *testing.T) {
	Check(io.EOF, WrapFunc(func(prev error) error {
		return nil
	}))
}

func TestCheckFuncMore(t *testing.T) {
	check := Check.With(Ignore(io.EOF))
	err := func() (err error) {
		defer Handle(&err)
		check(io.EOF)
		return
	}()
	if err != nil {
		t.Fatal()
	}

	err = func() (err error) {
		defer Handle(&err)
		check(nil, With(io.ErrNoProgress))
		return
	}()
	if err != nil {
		t.Fatal()
	}

	err = func() (err error) {
		defer Handle(&err)
		check(io.ErrClosedPipe)
		return
	}()
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}

	check = check.With(Ignore(io.ErrClosedPipe))
	err = func() (err error) {
		defer Handle(&err)
		check(io.EOF)
		return
	}()
	if err != nil {
		t.Fatal()
	}
	err = func() (err error) {
		defer Handle(&err)
		check(io.ErrClosedPipe)
		return
	}()
	if err != nil {
		t.Fatal()
	}
}
