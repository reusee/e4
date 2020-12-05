package e4

import (
	"fmt"
)

type Info string

var _ error = new(Info)

func (i *Info) Error() string {
	return string(*i)
}

func NewInfo(format string, args ...any) WrapFunc {
	return func(err error) error {
		i := Info(fmt.Sprintf(format, args...))
		return Chain{
			Err:  &i,
			Prev: err,
		}
	}
}

var WrapInfo = NewInfo

var WithInfo = NewInfo
