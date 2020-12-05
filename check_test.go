package e4

import (
	"io"
	"regexp"
	"testing"
)

func TestCheck(t *testing.T) {
	err := func() (err error) {
		defer Catch(&err)
		Check(io.EOF)
		return
	}()
	if !is(err, io.EOF) {
		t.Fatal()
	}

	err = func() (err error) {
		defer Catch(&err)
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

	err = func() (err error) {
		defer Catch(&err, NewInfo("foo %s", "bar"))
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
			defer Catch(nil)
			Check(io.EOF)
			return
		}()
	}()

	func() {
		defer func() {
			p := recover()
			if p != 42 {
				t.Fatal()
			}
		}()
		func() (err error) {
			defer Catch(&err)
			panic(42)
		}()
	}()

	err = func() (err error) {
		defer Catch(&err, NewInfo("foo %s", "bar"))
		return io.EOF
	}()
	if err.Error() != "foo bar\nEOF" {
		t.Fatal()
	}
	if !is(err, io.EOF) {
		t.Fatal()
	}

	err = func() (err error) {
		defer Catch(&err, NewInfo("foo %s", "bar"))
		return nil
	}()
	if err != nil {
		t.Fatal()
	}

	err = func() (err error) {
		defer Catch(&err)
		Check(io.EOF, func(err error) error {
			return nil
		})
		return
	}()
	if err != nil {
		t.Fatal()
	}

	err = func() (err error) {
		defer Catch(
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

}
