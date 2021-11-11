package e4

// CheckFunc is the type of Check, see Check's doc for details
type CheckFunc func(err error, args ...error) error

// Check is for error checking.
// if err is not nil, it will be wrapped by DefaultWrap then raised by Throw
var Check = CheckFunc(func(err error, args ...error) error {
	if err == nil {
		return nil
	}
	err = Wrap.With(args...)(err)
	return Throw(err)
})

// With returns a new CheckFunc that do additional wrapping with moreWraps
func (c CheckFunc) With(moreArgs ...error) CheckFunc {
	return func(err error, args ...error) error {
		if err == nil {
			return nil
		}
		err = Wrap.With(args...)(err)
		return c(err, moreArgs...)
	}
}

var CheckWithStacktrace = Check.With(WrapStacktrace)
