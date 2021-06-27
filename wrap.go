package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

// Wrap forms an error chain by calling wrap functions in order
func Wrap(err error, fns ...WrapFunc) error {
	if len(fns) == 0 {
		return err
	}
	fn := fns[0]
	fns = fns[1:]
	wrapped := fn(err)
	if wrapped == nil {
		return nil
	}
	if _, ok := wrapped.(Error); ok {
		return Wrap(wrapped, fns...)
	}
	return Wrap(MakeErr(wrapped, err), fns...)
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}
