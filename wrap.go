package e4

type WrapFunc func(err error) error

func Wrap(err error, fns ...WrapFunc) error {
	for _, fn := range fns {
		e := fn(err)
		if e != nil {
			if _, ok := e.(Error); !ok {
				err = MakeErr(e, err)
			} else {
				err = e
			}
		} else {
			return nil
		}
	}
	return err
}

func DefaultWrap(err error, fns ...WrapFunc) error {
	err = Wrap(err, fns...)
	if err != nil && !stacktraceIncluded(err) {
		err = NewStacktrace()(err)
	}
	return err
}
