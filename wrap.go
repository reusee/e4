package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

var Wrap = WrapFunc(func(err error) error {
	return err
})

func (w WrapFunc) With(fns ...WrapFunc) WrapFunc {
	return func(err error) error {
		for _, fn := range fns {
			if err == nil {
				return nil
			}
			wrapped := fn(err)
			if wrapped == nil {
				return nil
			}
			if _, ok := wrapped.(Error); !ok {
				wrapped = MakeErr(wrapped, err)
			}
			err = wrapped
		}
		return err
	}
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}
