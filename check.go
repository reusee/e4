package e4

// CheckFunc is the type of Check, see Check's doc for details
type CheckFunc func(err error, warpFuncs ...WrapFunc) error

// Check is for error checking.
// if err is not nil, it will be wrapped by DefaultWrap then raised by Throw
var Check = CheckFunc(func(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	return Throw(err)
})

// With returns a new CheckFunc that do additional wrapping with moreWraps
func (c CheckFunc) With(moreWraps ...WrapFunc) CheckFunc {
	return func(err error, fns ...WrapFunc) error {
		if err == nil {
			return nil
		}
		err = Wrap(err, fns...)
		return c(err, moreWraps...)
	}
}
