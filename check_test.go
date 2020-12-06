package e4

import (
	"io"
	"regexp"
	"testing"
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
		Check(io.EOF, NewInfo("foo %s", "bar"))
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
		defer Handle(&err, NewInfo("foo %s", "bar"))
		Check(io.EOF)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}
	ok, e = regexp.MatchString(
		"foo bar\n> at .*check_test.go:[0-9]+.*\n-.*\n-.*\n-.*\nEOF",
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
			te, ok := p.(error)
			if !ok {
				t.Fatal()
			}
			if !is(te, io.EOF) {
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
		defer Handle(&err, NewInfo("foo %s", "bar"))
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
		defer Handle(&err, NewInfo("foo %s", "bar"))
		return nil
	}()
	if err != nil {
		t.Fatal()
	}

	// set error to nil in check wrap func
	err = func() (err error) {
		defer Handle(&err)
		Check(io.EOF, func(err error) error {
			return nil
		})
		return
	}()
	if err != nil {
		t.Fatal()
	}

	// set error to nil in handle wrap func
	err = func() (err error) {
		defer Handle(
			&err,
			func(err error) error {
				return nil
			},
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
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}

}
