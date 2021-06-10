package e4

import (
	"fmt"
	"sync"
)

// Info represents a lazy-evaluaed formatted string
type Info struct {
	format     string
	str        string
	args       []any
	formatOnce sync.Once
}

var _ error = new(Info)

// Error implements error interface
func (i *Info) Error() string {
	i.formatOnce.Do(func() {
		i.str = fmt.Sprintf(i.format, i.args...)
	})
	return i.str
}

// NewInfo returns a WrapFunc that wraps an *Info error value
func NewInfo(format string, args ...any) WrapFunc {
	return With(&Info{
		format: format,
		args:   args,
	})
}
