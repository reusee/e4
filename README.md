# e4
error utilities version 4

### Example


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
	// handle error
	defer he(&err,
		// annotate error with string info
		e4.WithInfo("copy %s to %s", src, dst),
		// attach another error
		e4.With(fmt.Errorf("copy %s to %s", src, dst)),
		// attach a sentinel error value
		e4.With(ErrCopyFailed),
	)

	r, err := os.Open(src)
	// check error
	ce(err,
		// annotate error
		e4.WithInfo("open %s", src),
	)
	defer r.Close()

	w, err := os.Create(dst)
	ce(err,
		e4.WithInfo("create %s", dst),
		e4.With(ErrCreate{
			Path: dst,
		}),
	)
	// another error handling
	defer he(&err,
		// if error, close w
		e4.WithClose(w),
		// if error, remove dst file
		e4.WithFunc(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	// check error with if statement
	if err != nil {
		// throw error
		e4.Throw(err,
			e4.WithInfo("copy failed"),
		)
	}
	// check error
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

	err := CopyFile("demo.go", filepath.Join(os.TempDir(), "demo.go"))
	// calling Check without Handle will issue panic if error is returned
	ce(err)

	err = CopyFile("demo.go", "/")
	// check error with errors.Is / As
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
	  copy demo.go to /
	  $ main:demo.go:40 C:/Users/reus/reusee/e4/ main.CopyFile
	  & main:demo.go:86 C:/Users/reus/reusee/e4/ main.main
	  & runtime:proc.go:225 C:/Program Files/Go/src/runtime/ runtime.main
	  & runtime:asm_amd64.s:1371 C:/Program Files/Go/src/runtime/ runtime.goexit
	  create error: /
	  create /
	  open /: is a directory
	*/

}

```
