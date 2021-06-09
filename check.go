package e4

type CheckFunc func(err error, warpFuncs ...WrapFunc) error

var Check = CheckFunc(func(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	return Throw(err)
})

func CheckerWith(fns ...WrapFunc) func(error, ...WrapFunc) error {
	return func(err error, wrapFuncs ...WrapFunc) error {
		err = Wrap(err, fns...)
		return Check(err, wrapFuncs...)
	}
}

func (c CheckFunc) With(moreWraps ...WrapFunc) CheckFunc {
	return func(err error, fns ...WrapFunc) error {
		err = Wrap(err, fns...)
		return Check(err, moreWraps...)
	}
}
