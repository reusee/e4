package e4

import (
	"fmt"
	"sync"
)

// ErrInfo represents a lazy-evaluaed formatted string
type ErrInfo struct {
	format     string
	str        string
	args       []any
	formatOnce sync.Once
}

var _ error = new(ErrInfo)

// Error implements error interface
func (i *ErrInfo) Error() string {
	i.formatOnce.Do(func() {
		i.str = fmt.Sprintf(i.format, i.args...)
	})
	return i.str
}

// Info returns a WrapFunc that wraps an *ErrInfo error value
func Info(format string, args ...any) WrapFunc {
	return With(&ErrInfo{
		format: format,
		args:   args,
	})
}

var NewInfo = Info
