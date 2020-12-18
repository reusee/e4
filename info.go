package e4

import (
	"fmt"
	"sync"
)

type Info struct {
	format     string
	str        string
	args       []any
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
	return With(&Info{
		format: format,
		args:   args,
	})
}

var WrapInfo = NewInfo

var WithInfo = NewInfo
