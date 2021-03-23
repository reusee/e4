package e4

import "errors"

func Check(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	err = DefaultWrap(err, fns...)
	Throw(err)
}

func CheckerWith(fns ...WrapFunc) func(error, ...WrapFunc) {
	return func(err error, wrapFuncs ...WrapFunc) {
		err = Wrap(err, fns...)
		Check(err, wrapFuncs...)
	}
}

func DefaultWrap(err error, fns ...WrapFunc) error {
	err = Wrap(err, fns...)
	if err != nil {
		if !errors.Is(err, errStacktrace) {
			err = NewStacktrace()(err)
		}
	}
	return err
}

func Must(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	err = DefaultWrap(err, fns...)
	panic(err)
}

func Handle(errp *error, fns ...WrapFunc) {
	var err error
	if p := recover(); p != nil {
		if e, ok := p.(*throw); ok {
			err = e.err
		} else {
			panic(p)
		}
	}
	if errp != nil && *errp != nil {
		if err == nil {
			err = *errp
		} else {
			if !errors.Is(err, *errp) && !errors.Is(*errp, err) {
				err = MakeErr(err, *errp)
			}
		}
	}
	if err == nil {
		return
	}
	err = Wrap(err, fns...)
	if errp != nil {
		*errp = err
	} else {
		Throw(err)
	}
}
