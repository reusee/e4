package e4

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/reusee/e4/internal"
)

// Stacktrace represents call stack frames
type Stacktrace struct {
	Frames []Frame
}

// Frame represents a call frame
type Frame struct {
	File     string
	Dir      string
	Pkg      string
	Function string
	Line     int
	PkgPath  string
}

var _ error = new(Stacktrace)

// Error implements error interface
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

// NewStacktrace returns a WrapFunc that wraps current stacktrace
func NewStacktrace() WrapFunc {

	stacktrace := new(Stacktrace)
	v, put := pcsPool.Get()
	defer put()
	pcs := *(v.(*[]uintptr))
	skip := 1
	for {
		n := runtime.Callers(skip, pcs)
		if n == 0 {
			break
		}
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
			mod, fn := path.Split(frame.Function)
			if i := strings.Index(dir, mod); i > 0 {
				dir = dir[i:]
			}
			pkg := fn[:strings.IndexByte(fn, '.')]
			pkgPath := mod + pkg
			f := Frame{
				File:     file,
				Dir:      dir,
				Line:     frame.Line,
				Pkg:      pkg,
				Function: fn,
				PkgPath:  pkgPath,
			}
			stacktrace.Frames = append(stacktrace.Frames, f)
			if !more {
				break
			}
		}
		if n < len(pcs) {
			break
		}
	}
	return func(prev error) error {
		if prev == nil {
			return nil
		}
		if stacktraceIncluded(prev) {
			return prev
		}
		err := MakeErr(stacktrace, prev)
		err.flag |= flagStacktraceIncluded
		return err
	}
}

func stacktraceIncluded(err error) bool {
	if e, ok := err.(Error); ok &&
		e.flag&flagStacktraceIncluded > 0 {
		return true
	}
	if errors.As(err, new(*Stacktrace)) {
		return true
	}
	return false
}

var errStacktrace = errors.New("stacktrace")

// DropFrame returns a WrapFunc that drop Frames matching fn.
// If there is no existed stacktrace in chain, a new one will be created
func DropFrame(fn func(Frame) bool) WrapFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		var stacktrace *Stacktrace
		if !errors.As(err, &stacktrace) {
			err = NewStacktrace()(err)
			errors.As(err, &stacktrace)
		}
		newFrames := stacktrace.Frames[:0]
		for _, frame := range stacktrace.Frames {
			if fn(frame) {
				continue
			}
			newFrames = append(newFrames, frame)
		}
		stacktrace.Frames = newFrames
		return err
	}
}
