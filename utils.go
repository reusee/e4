package e4

import (
	"errors"
	"io"
	"strings"
)

func WithClose(c io.Closer) WrapFunc {
	return func(prev error) error {
		if err := c.Close(); err != nil {
			return With(err)(prev)
		}
		return prev
	}
}

func WithFunc(fn func()) WrapFunc {
	return func(prev error) error {
		fn()
		return prev
	}
}

func Ignore(err error) WrapFunc {
	return func(prev error) error {
		if errors.Is(prev, err) {
			return nil
		}
		return prev
	}
}

func IgnoreAs(target any) WrapFunc {
	return func(prev error) error {
		if errors.As(prev, target) {
			return nil
		}
		return prev
	}
}

func IgnoreContains(str string) WrapFunc {
	return func(prev error) error {
		if e := prev.Error(); strings.Contains(e, str) {
			return nil
		}
		return prev
	}
}
