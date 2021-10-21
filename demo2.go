//go:build ignore
// +build ignore

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
	wrap = e4.Wrap
)

func CopyFile(src, dst string) (err error) {
	wrap = wrap.With(
		e4.NewInfo("copy %s to %s", src, dst),
	)

	r, err := os.Open(src)
	if err != nil {
		return wrap(err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return wrap(err)
	}
	wrap = wrap.With(
		e4.Close(w),
		e4.Do(func() {
			os.Remove(dst)
		}),
	)

	_, err = io.Copy(w, r)
	if err != nil {
		return wrap(err)
	}

	if err := w.Close(); err != nil {
		return wrap(err)
	}

	return
}

func main() {

	err := CopyFile("demo.go", "/")

	var pathError *fs.PathError
	if !errors.As(err, &pathError) {
		panic("should be path error")
	}

	if err != nil {
		panic(err)
	}

}
