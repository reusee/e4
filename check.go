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
	err = Wrap(err, fns...)
	if err != nil {
		if trace := new(Stacktrace); !errors.As(err, &trace) {
			err = NewStacktrace()(err)
		}
	}
	panic(&check{
		err: err,
	})
}

func Must(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	err = Wrap(err, fns...)
	if err != nil {
		if trace := new(Stacktrace); !errors.As(err, &trace) {
			err = NewStacktrace()(err)
		}
	}
	panic(err)
}

func Handle(errp *error, fns ...WrapFunc) {
	var err error
	if errp != nil && *errp != nil {
		err = *errp
	}
	if p := recover(); p != nil {
		if e, ok := p.(*check); ok {
			if err != nil {
				err = Chain{
					Err:  e.err,
					Prev: err,
				}
			} else {
				err = e.err
			}
		} else {
			panic(p)
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
