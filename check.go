package e4

type CheckFunc func(err error, warpFuncs ...WrapFunc) error

var Check = CheckFunc(func(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	return Throw(err)
})

func (c CheckFunc) With(moreWraps ...WrapFunc) CheckFunc {
	return func(err error, fns ...WrapFunc) error {
		if err == nil {
			return nil
		}
		err = Wrap(err, fns...)
		return c(err, moreWraps...)
	}
}
