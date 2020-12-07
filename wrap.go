package e4

import "errors"

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

var _ error = WrapFunc(nil)

func (w WrapFunc) Error() string {
	return w(nil).Error()
}

func (w WrapFunc) Is(target error) bool {
	return errors.Is(w(nil), target)
}

func (w WrapFunc) As(target interface{}) bool {
	return errors.As(w(nil), target)
}
