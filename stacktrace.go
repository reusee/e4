package e4

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/reusee/e4/internal"
)

type Stacktrace struct {
	Frames []Frame
}

type Frame struct {
	File     string
	Dir      string
	Pkg      string
	Function string
	Line     int
}

var _ error = new(Stacktrace)

func (s *Stacktrace) Error() string {
	var b strings.Builder
	for i, frame := range s.Frames {
		if i == 0 {
			b.WriteString("$ ")
		} else {
			b.WriteString("\n& ")
		}
		b.WriteString(fmt.Sprintf(
			"%s:%s:%d %s %s",
			frame.Pkg,
			frame.File,
			frame.Line,
			frame.Dir,
			frame.Function,
		))
	}
	return b.String()
}

var pcsPool = internal.NewPool(
	128,
	func() any {
		bs := make([]uintptr, 32)
		return &bs
	},
)

func NewStacktrace() WrapFunc {
	stacktrace := new(Stacktrace)
	v, put := pcsPool.Get()
	defer put()
	pcs := *(v.(*[]uintptr))
	skip := 1
	for {
		n := runtime.Callers(skip, pcs)
		frames := runtime.CallersFrames(pcs[:n])
		for {
			skip++
			frame, more := frames.Next()
			if strings.HasPrefix(frame.Function, "github.com/reusee/e4.") &&
				!strings.HasPrefix(frame.Function, "github.com/reusee/e4.Test") {
				// internal funcs
				continue
			}
			dir, file := filepath.Split(frame.File)
			mod, fn := filepath.Split(frame.Function)
			if i := strings.Index(dir, mod); i > 0 {
				dir = dir[i:]
			}
			var pkg string
			if fn != "" {
				pkg = fn[:strings.IndexByte(fn, '.')]
			}
			stacktrace.Frames = append(stacktrace.Frames, Frame{
				File:     file,
				Dir:      dir,
				Line:     frame.Line,
				Pkg:      pkg,
				Function: fn,
			})
			if !more {
				break
			}
		}
		if n < len(pcs) {
			break
		}
	}
	return func(prev error) error {
		err := MakeErr(stacktrace, prev)
		err.flag |= flagStacktraceIncluded
		return err
	}
}

func stacktraceIncluded(err error) bool {
	if e, ok := err.(Error); !ok {
		return false
	} else {
		return e.flag&flagStacktraceIncluded > 0
	}
}

var errStacktrace = errors.New("stacktrace")

var WithStacktrace = NewStacktrace

var WrapStacktrace = NewStacktrace

var WithStack = NewStacktrace

var WrapStack = NewStacktrace
