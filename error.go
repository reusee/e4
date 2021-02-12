package e4

import (
	"errors"
	"strings"
)

type Error struct {
	Err    error
	Prev   error
	Bubble error
}

func (c Error) Is(target error) bool {
	if c.Bubble != nil && errors.Is(c.Bubble, target) {
		return true
	}
	return errors.Is(c.Err, target)
}

func (c Error) As(target interface{}) bool {
	if c.Bubble != nil && errors.As(c.Bubble, target) {
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
		Bubble: getBubble(prev),
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
		return e.Bubble
	}
}
