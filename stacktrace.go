package e4

import (
	"fmt"
	"runtime"
	"strings"
)

type Stacktrace struct {
	Frames []Frame
}

type Frame struct {
	File     string
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
		b.WriteString(fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	}
	return b.String()
}

func NewStacktrace() WrapFunc {
	return func(err error) error {
		stacktrace := new(Stacktrace)
		numPCs := 32
		for {
			pcs := make([]uintptr, numPCs)
			n := runtime.Callers(1, pcs)
			if n == len(pcs) { // NOCOVER
				numPCs *= 2
				continue
			}
			pcs = pcs[:n]
			frames := runtime.CallersFrames(pcs)
			for {
				frame, more := frames.Next()
				if strings.HasPrefix(frame.Function, "github.com/reusee/e4.") &&
					!strings.HasPrefix(frame.Function, "github.com/reusee/e4.Test") {
					// internal funcs
					continue
				}
				stacktrace.Frames = append(stacktrace.Frames, Frame{
					File:     frame.File,
					Line:     frame.Line,
					Function: frame.Function,
				})
				if !more {
					break
				}
			}
			break
		}
		return Chain{
			Err:  stacktrace,
			Prev: err,
		}
	}
}

var WithStacktrace = NewStacktrace

var WrapStacktrace = NewStacktrace

var WithStack = NewStacktrace

var WrapStack = NewStacktrace
