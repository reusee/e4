# e4
Error handling utilities

## Features

* Ad-hoc error wrapping
  + easy to chain multiple error values
  + no need to implement Unwrap / Is / As for every error type
  + inspect error hierarchy with standard errors.Is / As
* Alternatives to `if err != nil { return ... }` statement
  + utilizing panic / recover
  + but not crossing the function boundary
  + not forced to use, play well with existing codes

## CopyFile demo

```go
package main

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/reusee/e4"
)

var (
	// ergonomic aliases
	check, handle = e4.Check, e4.Handle
)

func CopyFile(src, dst string) (err error) {
	defer handle(&err,
		e4.NewInfo("copy %s to %s", src, dst),
	)

	r, err := os.Open(src)
	check(err)
	defer r.Close()

	w, err := os.Create(dst)
	check(err)
	defer handle(&err,
		e4.Close(w),
		e4.Do(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	check(err)

	check(w.Close())

	return
}

func main() {

	err := CopyFile("demo.go", "/")

	var pathError *fs.PathError
	if !errors.As(err, &pathError) {
		panic("should be path error")
	}

	check(err)

}
```

## Usages

### Error wrapping

errors can be wrapped with `e4.Wrap` and various util functions.

```go
err := foo()

if err != nil {
  return e4.Wrap.With(

    // wrap another error value
    e4.With(io.ErrUnexpectedEOF),

    // wrap a lazy-formatted message
    e4.NewInfo("unexpected %s", "EOF"),

    // close a io.Closer
    e4.Close(w),

    // do something
    e4.Do(func() {
      fmt.Printf("error occur\n")
    }),

    // cumstom wrap function
    func(err error) error {
      fmt.Printf("wrap EOF\n")
      return e4.MakeErr(err, io.EOF)
    },

  )(err)
}
```

wrapped errors can be inspected with `errors.Is` or `errors.As`

```go
errFoo := errors.New("foo")

err := e4.Wrap.With(
  e4.With(io.ErrUnexpectedEOF),
  e4.With(errFoo),
  e4.With(new(fs.PathError)),
  // wrap a nested error
  e4.With(e4.Wrap.With(fs.ErrInvalid,
    e4.With(e4.Wrap.With(io.ErrClosedPipe,
      e4.With(io.ErrShortWrite))))),
)(io.EOF)

errors.Is(err, io.EOF) // true
errors.Is(err, io.ErrUnexpectedEOF) // true
errors.Is(err, io.ErrShortWrite) // true for deeply nested values
var pathError *fs.PathError
errors.As(err, &pathError) // true
```

### Alternatives to if and return statement 

error values can be thrown and catch

```go
func foo() (err error) {
  defer e4.Handle(&err)
  if err := bar(); err != nil {
    e4.Throw(err)
  }
}
```

Further, `if` and `e4.Throw` can be replaced with `e4.Check`

```go
func foo() (err error) {
  defer e4.Handle(&err)
  e4.Check(bar())
}
```

Error wrapping also works in the Check site or the Handle site

```go
func foo() (err error) {
  defer e4.Handle(&err,
    e4.With(fmt.Errorf("foo error")),
    e4.Do(func() {
      fmt.Printf("foo error\n")
    }),
  )
  e4.Check(bar(),
    e4.With(fmt.Errorf("bar error")),
    // wrap stack trace
    e4.WrapStacktrace,
    // ignore errors that errors.Is return true
    e4.Ignore(io.EOF),
    // ignore errors that errors.As return true
    e4.IgnoreAs(new(*fs.PathError)),
    // ignore errors that Error() contains specific string
    e4.IgnoreContains("EOF"),
  )
}
```

### Check with stacktrace

Errors checked by `e4.CheckWithStacktrace` are implicitly wrapped by stacktrace.

```go
func foo() (err error) {
  defer e4.Handle(&err)
  e4.CheckWithStacktrace(io.EOF)
}

err := foo()

var trace *e4.Stacktrace
errors.As(err, &trace) // true

```

To drop not-interested frames from the stacktrace, decorate the `e4.CheckWithStacktrace` with `e4.DropFrame`

```go
var check = e4.CheckWithStacktrace.With(e4.DropFrame(func(frame e4.Frame) bool {
  // drop runtime and reflect frames
  return frame.PkgPath == "runtime" || 
    frame.PkgPath == "reflect"
}))

check(io.EOF)
```

