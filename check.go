package e4

import "errors"

type thrownError struct {
	err error
	sig int64
}

func (t *thrownError) String() string { // NOCOVER
	return t.err.Error()
}

func Check(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	for _, fn := range fns {
		e := fn(err)
		if e != nil {
			if _, ok := e.(Chain); !ok {
				err = Chain{
					Err:  e,
					Prev: err,
				}
			} else {
				err = e
			}
		} else {
			err = e
		}
	}
	if err != nil {
		if trace := new(Stacktrace); !errors.As(err, &trace) {
			err = NewStacktrace()(err)
		}
	}
	panic(&thrownError{
		err: err,
	})
}

func Handle(errp *error, fns ...WrapFunc) {
	var err error
	if errp != nil && *errp != nil {
		err = *errp
	}
	if p := recover(); p != nil {
		if e, ok := p.(*thrownError); ok {
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
	for _, fn := range fns {
		e := fn(err)
		if e != nil {
			if _, ok := e.(Chain); !ok {
				err = Chain{
					Err:  e,
					Prev: err,
				}
			} else {
				err = e
			}
		} else {
			err = e
		}
	}
	if errp != nil {
		*errp = err
	} else {
		panic(err)
	}
}
