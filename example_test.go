package e4

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func ExampleCheck_handle() {

	fn := func() (err error) {
		defer Handle(&err)
		Check(io.EOF)
		return
	}

	err := fn()
	fmt.Println(errors.Is(err, io.EOF))

	// Check wraps stack trace automatically
	var stacktrace *Stacktrace
	if errors.As(err, &stacktrace) {
		fmt.Println(len(stacktrace.Frames) > 0)
	}

	// Output:
	// true
	// true
}

func ExampleCheck_wrap_function() {

	wrapWithChain := func(err error) error {
		return MakeErr(os.ErrClosed, err)
	}

	wrapAutoChain := func(err error) error {
		return os.ErrDeadlineExceeded
	}

	fn := func() (err error) {
		defer Handle(&err)
		Check(
			io.EOF,
			wrapWithChain,
			wrapAutoChain,
		)
		return
	}

	err := fn()
	fmt.Println(errors.Is(err, io.EOF))
	fmt.Println(errors.Is(err, os.ErrClosed))
	fmt.Println(errors.Is(err, os.ErrDeadlineExceeded))

	// Output:
	// true
	// true
	// true
}

func ExampleHandle_wrap_function() {

	wrapWithChain := func(err error) error {
		return MakeErr(os.ErrClosed, err)
	}

	wrapAutoChain := func(err error) error {
		return os.ErrDeadlineExceeded
	}

	fn := func() (err error) {
		defer Handle(
			&err,
			wrapWithChain,
			wrapAutoChain,
		)
		Check(io.EOF)
		return
	}

	err := fn()
	fmt.Println(errors.Is(err, io.EOF))
	fmt.Println(errors.Is(err, os.ErrClosed))
	fmt.Println(errors.Is(err, os.ErrDeadlineExceeded))

	// Output:
	// true
	// true
	// true
}

func ExampleHandle_wrap_return_error() {

	fn := func(name string) (err error) {
		defer Handle(&err, WithInfo("hello, %s", name))
		return io.EOF
	}

	err := fn("world")
	fmt.Println(errors.Is(err, io.EOF))
	var errInfo *Info
	if errors.As(err, &errInfo) {
		fmt.Println(errInfo.Error())
	}

	// Output:
	// true
	// hello, world
}

func ExampleHandle_igonre() {

	ignoreEOF := func(err error) error {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	fn := func() (err error) {
		defer Handle(&err, ignoreEOF)
		return io.EOF
	}

	err := fn()
	fmt.Println(err)

	// Output:
	// <nil>
}

func ExampleCheck_copyFile() {

	copyFile := func(src, dst string) (err error) {
		defer Handle(&err,
			WithInfo("copy %s to %s", src, dst),
		)

		r, err := os.Open(src)
		Check(err, WithInfo("open %s", src))
		defer r.Close()

		w, err := os.Create(dst)
		Check(err, WithInfo("create %s", dst))
		defer w.Close()
		defer Handle(&err, WithFunc(func() {
			os.Remove(dst)
		}))

		_, err = io.Copy(w, r)
		Check(err)
		Check(w.Close())

		return
	}

	Check(copyFile("example_test.go", "example_test.go"))

}
