package e4

import (
	"errors"
	"io"
	"testing"
)

var testErr error

func BenchmarkHandleNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Handle(&err)
			return nil
		}()
		if err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkHandleErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Handle(&err)
			return io.EOF
		}()
		if !errors.Is(err, io.EOF) {
			b.Fatal()
		}
	}
}

func BenchmarkHandleCheckNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Handle(&err)
			Check(nil)
			return
		}()
		if err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkHandleCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := func() (err error) {
			defer Handle(&err)
			Check(io.EOF)
			return
		}()
		if !errors.Is(err, io.EOF) {
			b.Fatal()
		}
	}
}

func BenchmarkHandleCheckDeep(b *testing.B) {
	var fn func(int) error
	fn = func(i int) (err error) {
		defer Handle(&err)
		if i == 0 {
			return io.EOF
		}
		Check(fn(i - 1))
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := fn(30)
		if !errors.Is(err, io.EOF) {
			b.Fatal()
		}
	}
}
