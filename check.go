package e4

import "errors"

type check struct {
	err error
}

func (c *check) String() string { // NOCOVER
	return c.err.Error()
}

func Check(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	err = DefaultWrap(err, fns...)
	panic(&check{
		err: err,
	})
}

func DefaultWrap(err error, fns ...WrapFunc) error {
	err = Wrap(err, fns...)
	if err != nil {
		if trace := new(Stacktrace); !errors.As(err, &trace) {
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
		if e, ok := p.(*check); ok {
			err = e.err
		} else {
			panic(p)
		}
	}
	// consider **errp only when err is nil
	if err == nil {
		if errp != nil && *errp != nil {
			err = *errp
		}
	}
	if err == nil {
		return
	}
	err = Wrap(err, fns...)
	if errp != nil {
		*errp = err
	} else {
		panic(err)
	}
}
