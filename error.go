package e4

import (
	"errors"
	"strings"
)

// Error represents a chain of errors
type Error struct {
	Err  error
	Prev error

	flag uint64
}

// Is reports whether any error in the chain matches target
func (c Error) Is(target error) bool {
	return errors.Is(c.Err, target)
}

// As reports whether any error in the chain matches target.
// And if so, assign the first matching error to target
func (c Error) As(target interface{}) bool {
	return errors.As(c.Err, target)
}

// Unwrap returns Prev error
func (c Error) Unwrap() error {
	return c.Prev
}

// Error implements error interface
func (c Error) Error() string {
	var b strings.Builder
	b.WriteString(c.Err.Error())
	if c.Prev != nil {
		b.WriteString("\n")
		b.WriteString(c.Prev.Error())
	}
	return b.String()
}

// MakeErr creates an Error with internal manipulations.
// It's safe to construct Error without calling this function but not encouraged
func MakeErr(err error, prev error) Error {
	return Error{
		Err:  err,
		Prev: prev,
		flag: getFlag(err) | getFlag(prev),
	}
}

// With returns a WrapFunc that wraps an error value
func With(err error) WrapFunc {
	return func(prev error) error {
		if prev == nil {
			return nil
		}
		return MakeErr(err, prev)
	}
}

func getFlag(err error) uint64 {
	if e, ok := err.(Error); !ok {
		return 0
	} else {
		return e.flag
	}
}
