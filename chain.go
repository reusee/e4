package e4

import (
	"errors"
	"strings"
)

type Chain struct {
	Err  error
	Prev error
}

func (c Chain) Is(target error) bool {
	return errors.Is(c.Err, target)
}

func (c Chain) As(target interface{}) bool {
	return errors.As(c.Err, target)
}

func (c Chain) Unwrap() error {
	return c.Prev
}

func (c Chain) Error() string {
	var b strings.Builder
	b.WriteString(c.Err.Error())
	if c.Prev != nil {
		b.WriteString("\n")
		b.WriteString(c.Prev.Error())
	}
	return b.String()
}

func With(err error) WrapFunc {
	return func(prev error) error {
		return Chain{
			Err:  err,
			Prev: prev,
		}
	}
}
