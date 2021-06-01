// +build ignore

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
		// chain more error values
		e4.WithInfo("copy %s to %s", src, dst),
		e4.With(fmt.Errorf("copy %s to %s", src, dst)),
		e4.With(ErrCopyFailed),
	)

	r, err := os.Open(src)
	// check error
	ce(err,
		// chain info error
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
		e4.Close(w),
		// if error, remove dst file
		e4.Do(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	// check error with if statement and Throw
	if err != nil {
		e4.Throw(err,
			e4.WithInfo("copy failed"),
		)
	}

	// check error with if statement and Wrap
	if err := w.Close(); err != nil {
		return e4.Wrap(err,
			e4.WithInfo("close failed"),
		)
	}

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
	// match errors in chain with Is / As
	if !errors.Is(err, ErrCopyFailed) {
		panic("shoule be ErrCopyFailed")
	}
	if !errors.As(err, new(ErrCreate)) {
		panic("should be ErrCreate")
	}

	// stacktrace is added automatically
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
