package e4

func Check(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	return Throw(err)
}

func CheckerWith(fns ...WrapFunc) func(error, ...WrapFunc) error {
	return func(err error, wrapFuncs ...WrapFunc) error {
		err = Wrap(err, fns...)
		return Check(err, wrapFuncs...)
	}
}
