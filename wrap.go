package e4

type WrapFunc func(err error) error

func Wrap(err error, fns ...WrapFunc) error {
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
	return err
}
