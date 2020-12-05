package e4

import (
	"fmt"
	"sync"
)

type Info struct {
	format     string
	args       []any
	str        string
	formatOnce sync.Once
}

var _ error = new(Info)

func (i *Info) Error() string {
	i.formatOnce.Do(func() {
		i.str = fmt.Sprintf(i.format, i.args...)
	})
	return i.str
}

func NewInfo(format string, args ...any) WrapFunc {
	return func(err error) error {
		return Chain{
			Err: &Info{
				format: format,
				args:   args,
			},
			Prev: err,
		}
	}
}

var WrapInfo = NewInfo

var WithInfo = NewInfo
