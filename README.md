# e4
Error handling utilities

## Features

* Ad-hoc error wrapping
* Auto stacktrace wrapping
* Alternatives to `if err != nil { return ... }` statement

## Usages

### Error wrapping

errors can be wrapped with `e4.Wrap` and various util functions.

```go
err := foo()

if err != nil {
  return e4.Wrap(err,

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

  )
}
```

wrapped errors can be inspected with `errors.Is` or `errors.As`

```go
errFoo := errors.New("foo")

err := e4.Wrap(io.EOF,
  e4.With(io.ErrUnexpectedEOF),
  e4.With(errFoo),
  e4.With(new(fs.PathError)),
  // wrap a nested error
  e4.With(e4.Wrap(fs.ErrInvalid,
    e4.With(e4.Wrap(io.ErrClosedPipe,
      e4.With(io.ErrShortWrite))))),
)

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
    // ignore errors that errors.Is return true
    e4.Ignore(io.EOF),
    // ignore errors that errors.As return true
    e4.IgnoreAs(new(*fs.PathError)),
    // ignore errors that Error() contains specific string
    e4.IgnoreContains("EOF"),
  )
}
```

### Auto wrapping stacktrace

Errors checked by `e4.Check` are implicitly wrapped by stacktrace.

```go
func foo() (err error) {
  defer e4.Handle(&err)
  e4.Check(io.EOF)
}

err := foo()

var trace *e4.Stacktrace
errors.As(err, &trace) // true

```

The `e4.DefaultWrap` function also wraps stacktrace automatically.

To drop not-interested frames from the stacktrace, decorate the `e4.Check` with `e4.DropFrame`

```go
var check = e4.Check.With(e4.DropFrame(func(frame e4.Frame) bool {
  // drop runtime and reflect frames
  return frame.PkgPath == "runtime" || 
    frame.PkgPath == "reflect"
}))

check(io.EOF)
```

### A file copy demo

```go
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/reusee/e4"
)

var (
	// ergonomic aliases
	ce, he = e4.Check, e4.Handle
)

func CopyFile(src, dst string) (err error) {
	defer he(&err,
		e4.NewInfo("copy %s to %s", src, dst),
		e4.With(ErrCopyFailed),
	)

	r, err := os.Open(src)
	ce(err,
		e4.NewInfo("open %s", src),
	)
	defer r.Close()

	w, err := os.Create(dst)
	ce(err,
		e4.NewInfo("create %s", dst),
		e4.With(ErrCreate{
			Path: dst,
		}),
	)
	defer he(&err,
		e4.Close(w),
		e4.Do(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	ce(err)

	ce(w.Close())

	return
}

var ErrCopyFailed = errors.New("copy failed")

type ErrCreate struct {
	Path string
}

func (e ErrCreate) Error() string {
	return fmt.Sprintf("create error: %s", e.Path)
}

func main() {

	err := CopyFile(
		"demo.go",
		filepath.Join(os.TempDir(), "demo.go"),
	)
	ce(err)

	err = CopyFile("demo.go", "/")
	if !errors.Is(err, ErrCopyFailed) {
		panic("shoule be ErrCopyFailed")
	}
	if !errors.As(err, new(ErrCreate)) {
		panic("should be ErrCreate")
	}

	println(err.Error())
	/*
	   copy failed
	   copy demo.go to /
	   $ main:demo.go:33 C:/Users/reus/reusee/e4/ main.CopyFile
	   & main:demo.go:72 C:/Users/reus/reusee/e4/ main.main
	   & runtime:proc.go:225 C:/Program Files/Go/src/runtime/ runtime.main
	   & runtime:asm_amd64.s:1371 C:/Program Files/Go/src/runtime/ runtime.goexit
	   create error: /
	   create /
	   open /: is a directory
	*/

}
```
