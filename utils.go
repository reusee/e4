package e4

import (
	"errors"
	"io"
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
