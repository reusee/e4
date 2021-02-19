# e4
error utilities version 4

### Example


```
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
		// chain another string error
		e4.With(fmt.Errorf("copy %s to %s", src, dst)),
		// chain a sentinel error value
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
		e4.Throw(err)
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

}
```
