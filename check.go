package e4

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
		err = fn(err)
	}
	if err != nil {
		if trace := new(Stacktrace); !as(err, &trace) {
			err = NewStacktrace()(err)
		}
	}
	panic(&thrownError{
		err: err,
	})
}

func Catch(errp *error, fns ...WrapFunc) {
	var err error
	if p := recover(); p != nil {
		if e, ok := p.(*thrownError); ok {
			err = e.err
		} else {
			panic(p)
		}
	} else {
		if errp != nil && *errp != nil {
			err = *errp
		}
	}
	if err == nil {
		return
	}
	for _, fn := range fns {
		err = fn(err)
	}
	if errp != nil {
		*errp = err
	} else {
		panic(err)
	}
}

var C = Check

var T = Catch

var Handle = Catch
