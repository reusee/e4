package e4

import (
	"errors"
	"strings"
)

type Error struct {
	Err  error
	Prev error

	bubble error
}

func (c Error) Is(target error) bool {
	if c.bubble != nil && errors.Is(c.bubble, target) {
		return true
	}
	return errors.Is(c.Err, target)
}

func (c Error) As(target interface{}) bool {
	if c.bubble != nil && errors.As(c.bubble, target) {
		return true
	}
	return errors.As(c.Err, target)
}

func (c Error) Unwrap() error {
	return c.Prev
}

func (c Error) Error() string {
	var b strings.Builder
	b.WriteString(c.Err.Error())
	if c.Prev != nil {
		b.WriteString("\n")
		b.WriteString(c.Prev.Error())
	}
	return b.String()
}

func MakeErr(err error, prev error) Error {
	return Error{
		Err:    err,
		Prev:   prev,
		bubble: getBubble(prev),
	}
}

func With(err error) WrapFunc {
	return func(prev error) error {
		return MakeErr(err, prev)
	}
}

func getBubble(err error) error {
	if e, ok := err.(Error); !ok {
		return nil
	} else {
		return e.bubble
	}
}
