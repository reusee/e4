package e4

type Debug struct {
	*Info
}

var _ error = new(Debug)

func (d Debug) ErrorLevel() Level {
	return DebugLevel
}

func NewDebug(format string, args ...any) WrapFunc {
	return With(&Debug{
		Info: &Info{
			format: format,
			args:   args,
		},
	})
}
