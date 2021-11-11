package e4

import "testing"

// WrapFunc wraps an error to form a chain.
//
// Instances must follow these rules:
// if argument is nil, return value must be nil
type WrapFunc func(err error) error

var _ error = WrapFunc(nil)

var Wrap = WrapFunc(func(err error) error {
	return err
})

func (w WrapFunc) With(args ...error) WrapFunc {
	// convert to WrapFuncs
	for i, arg := range args {
		switch arg := arg.(type) {
		case WrapFunc:
		case error:
			args[i] = With(arg)
		}
	}
	return func(err error) error {
		for _, arg := range args {
			if err == nil {
				return nil
			}
			wrapped := arg.(WrapFunc)(err)
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

func (w WrapFunc) Error() string {
	panic("should not used as error")
}

func (w WrapFunc) Assign(p *error) {
	if p != nil && *p != nil {
		*p = w(*p)
	}
}

// TestWrapFunc tests a WrapFunc instance
func TestWrapFunc(t *testing.T, fn WrapFunc) {
	if fn(nil) != nil {
		t.Fatal("should return nil")
	}
}
