package e4

import (
	"errors"
	"io"
	"strings"
)

// Close returns a WrapFunc that closes the Closer
func Close(c io.Closer) WrapFunc {
	return func(prev error) error {
		if err := c.Close(); err != nil {
			return With(err)(prev)
		}
		return prev
	}
}

// Do returns a WrapFunc that calls fn
func Do(fn func()) WrapFunc {
	return func(prev error) error {
		fn()
		return prev
	}
}

// Ignore returns a WrapFunc that returns nil if errors.Is(prev, err) is true
func Ignore(err error) WrapFunc {
	return func(prev error) error {
		if errors.Is(prev, err) {
			return nil
		}
		return prev
	}
}

// IgnoreAs returns a WrapFunc that returns nil if errors.As(prev, target) is true
func IgnoreAs(target any) WrapFunc {
	return func(prev error) error {
		if errors.As(prev, target) {
			return nil
		}
		return prev
	}
}

// IgnoreContains returns a WrapFunc that returns nil if prev.Error() contains str
func IgnoreContains(str string) WrapFunc {
	return func(prev error) error {
		if prev == nil {
			return nil
		}
		if e := prev.Error(); strings.Contains(e, str) {
			return nil
		}
		return prev
	}
}
