package e4

import (
	"errors"
	"strings"
)

type Error struct {
	Err  error
	Prev error

	flag uint64
}

func (c Error) Is(target error) bool {
	return errors.Is(c.Err, target)
}

func (c Error) As(target interface{}) bool {
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
		Err:  err,
		Prev: prev,
		flag: getFlag(err) | getFlag(prev),
	}
}

func With(err error) WrapFunc {
	return func(prev error) error {
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
