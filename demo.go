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
