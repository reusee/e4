package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

// Wrap combines multiple WrapFuncs sequentially
func Wrap(fns ...WrapFunc) WrapFunc {
	return func(err error) error {
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
			return Wrap(fns[1:]...)(wrapped)
		}
		return Wrap(fns[1:]...)(MakeErr(wrapped, err))
	}
}

func (w WrapFunc) With(fns ...WrapFunc) WrapFunc {
	return func(err error) error {
		return Wrap(fns...)(w(err))
	}
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}

var WrapWithStacktrace = Wrap(WrapStacktrace)
