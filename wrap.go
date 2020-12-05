package e4

type WrapFunc func(err error) error

func Wrap(err error, fns ...WrapFunc) error {
	for _, fn := range fns {
		err = fn(err)
	}
	return err
}

var _ error = WrapFunc(nil)

func (w WrapFunc) Error() string {
	return w(nil).Error()
}

var W = Wrap
