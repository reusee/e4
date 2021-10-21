package e4

type ErrDebug struct {
	*ErrInfo
}

var _ error = new(ErrDebug)

func (d ErrDebug) ErrorLevel() Level {
	return DebugLevel
}

func Debug(format string, args ...any) WrapFunc {
	return With(&ErrDebug{
		ErrInfo: &ErrInfo{
			format: format,
			args:   args,
		},
	})
}

var NewDebug = Debug
