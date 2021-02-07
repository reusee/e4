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
	defer he(&err,
		e4.WithInfo("copy %s to %s", src, dst),
		e4.With(fmt.Errorf("copy %s to %s", src, dst)),
		e4.With(ErrCopyFailed),
	)

	r, err := os.Open(src)
	ce(err,
		e4.WithInfo("open %s", src),
	)
	defer r.Close()

	w, err := os.Create(dst)
	ce(err,
		e4.WithInfo("create %s", dst),
	)
	defer he(&err,
		e4.WithClose(w),
		e4.WithFunc(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	ce(err)
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

}
