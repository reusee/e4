package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

// Wrap forms an error chain by calling wrap functions in order
func Wrap(err error, fns ...WrapFunc) error {
	for _, fn := range fns {
		e := fn(err)
		if e != nil {
			if _, ok := e.(Error); !ok {
				err = MakeErr(e, err)
			} else {
				err = e
			}
		} else {
			return nil
		}
	}
	return err
}

// DefaultWrap wraps error with stacktrace
func DefaultWrap(err error, fns ...WrapFunc) error {
	err = Wrap(err, fns...)
	if err != nil && !stacktraceIncluded(err) {
		err = NewStacktrace()(err)
	}
	return err
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}
