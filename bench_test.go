package e4

import (
	"errors"
	"io"
	"testing"
)

var testErr error

func BenchmarkCatchNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Catch(&err)
			return nil
		}()
		if err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkCatchErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Catch(&err)
			return io.EOF
		}()
		if !errors.Is(err, io.EOF) {
			b.Fatal()
		}
	}
}

func BenchmarkCatchCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Catch(&err)
			Check(io.EOF)
			return
		}()
		if !errors.Is(err, io.EOF) {
			b.Fatal()
		}
	}
}

func BenchmarkCatchCheckNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Catch(&err)
			Check(nil)
			return
		}()
		if err != nil {
			b.Fatal()
		}
	}
}
