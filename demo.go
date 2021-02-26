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
	ce, he = e4.Check, e4.Handle
)

func CopyFile(src, dst string) (err error) {
	// handle/catch error
	defer he(&err,
		// annotate error with string info
		e4.WithInfo("copy %s to %s", src, dst),
		// attach another string error
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
		e4.Throw(
			// wrap errors manually
			e4.Wrap(err,
				e4.WithInfo("copy failed"),
			),
		)
	}
	ce(w.Close())

	return
}

var ErrCopyFailed = errors.New("copy failed")

func main() {

	err := CopyFile("demo.go", filepath.Join(os.TempDir(), "demo.go"))
	ce(err)

	err = CopyFile("demo.go", "/")
	if !errors.Is(err, ErrCopyFailed) {
		panic("shoule be ErrCopyFailed")
	}

	println(err.Error())
	/*
	  copy demo.go to /
	  copy demo.go to /
	  $ main.demo.go:39 /home/reus/reusee/e4/ main.CopyFile
	  & main.demo.go:75 /home/reus/reusee/e4/ main.main
	  & runtime.proc.go:225 /usr/lib/go/src/runtime/ runtime.main
	  & runtime.asm_amd64.s:1371 /usr/lib/go/src/runtime/ runtime.goexit
	  create /
	  open /: is a directory
	*/

}
