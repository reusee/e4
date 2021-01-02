package e4

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Stacktrace struct {
	Frames []Frame
}

type Frame struct {
	File     string
	Dir      string
	Function string
	Line     int
}

var _ error = new(Stacktrace)

func (s *Stacktrace) Error() string {
	var b strings.Builder
	for i, frame := range s.Frames {
		if i == 0 {
			b.WriteString("> at ")
		} else {
			b.WriteString("\n-    ")
		}
		//b.WriteString(fmt.Sprintf("%s:%d ## %s", frame.File, frame.Line, frame.Function))
		b.WriteString(fmt.Sprintf(
			"%s:%d %s %s",
			frame.File,
			frame.Line,
			frame.Dir,
			frame.Function,
		))
	}
	return b.String()
}

func NewStacktrace() WrapFunc {
	stacktrace := new(Stacktrace)
	pcs := make([]uintptr, 32)
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
			stacktrace.Frames = append(stacktrace.Frames, Frame{
				File:     file,
				Dir:      dir,
				Line:     frame.Line,
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
	return With(stacktrace)
}

var WithStacktrace = NewStacktrace

var WrapStacktrace = NewStacktrace

var WithStack = NewStacktrace

var WrapStack = NewStacktrace
