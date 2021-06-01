package e4

import "errors"

func Check(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	return Throw(err)
}

func CheckPtr(ptr *error, fns ...WrapFunc) error {
	if ptr == nil {
		return nil
	}
	if *ptr == nil {
		return nil
	}
	err := DefaultWrap(*ptr, fns...)
	*ptr = err
	return Throw(err)
}

func CheckerWith(fns ...WrapFunc) func(error, ...WrapFunc) error {
	return func(err error, wrapFuncs ...WrapFunc) error {
		err = Wrap(err, fns...)
		return Check(err, wrapFuncs...)
	}
}

func PtrCheckerWith(fns ...WrapFunc) func(*error, ...WrapFunc) error {
	return func(ptr *error, wrapFuncs ...WrapFunc) error {
		if ptr == nil {
			return nil
		}
		*ptr = Wrap(*ptr, fns...)
		return CheckPtr(ptr, wrapFuncs...)
	}
}

func DefaultWrap(err error, fns ...WrapFunc) error {
	err = Wrap(err, fns...)
	if err != nil && !stacktraceIncluded(err) {
		err = NewStacktrace()(err)
	}
	return err
}

func Must(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	panic(err)
}

func Handle(errp *error, fns ...WrapFunc) {
	var err error
	// check throw error
	if p := recover(); p != nil {
		if e, ok := p.(*throw); ok {
			err = e.err
		} else {
			panic(p)
		}
	}
	// check pointed error
	if errp != nil && *errp != nil {
		if err == nil {
			// no throw error
			err = *errp
		} else {
			// wrap if not the same
			if !errors.Is(err, *errp) && !errors.Is(*errp, err) {
				err = MakeErr(err, *errp)
			}
		}
	}
	if err == nil {
		// no error
		return
	}
	// wrap
	err = Wrap(err, fns...)
	if errp != nil {
		// set pointed variable
		*errp = err
	} else {
		// throw
		Throw(err)
	}
}
