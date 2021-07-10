package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

// Wrap forms an error chain by calling wrap functions in order
func Wrap(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	if len(fns) == 0 {
		return err
	}
	wrapped := fns[0](err)
	if wrapped == nil {
		return nil
	}
	if _, ok := wrapped.(Error); ok {
		return Wrap(wrapped, fns[1:]...)
	}
	return Wrap(MakeErr(wrapped, err), fns[1:]...)
}

func (w WrapFunc) With(fns ...WrapFunc) WrapFunc {
	return func(err error) error {
		return Wrap(w(err), fns...)
	}
}

func (w WrapFunc) Wrap(err error, fns ...WrapFunc) error {
	return Wrap(w(err), fns...)
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}

var WrapWithStacktrace = WrapStacktrace.Wrap
