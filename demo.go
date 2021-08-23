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
