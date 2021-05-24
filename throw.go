package e4

type throw struct {
	err error
}

func (c *throw) String() string { // NOCOVER
	return c.err.Error()
}

func (c *throw) Error() string { // NOCOVER
	return c.err.Error()
}

func (c *throw) Unwrap() error {
	return c.err
}

func Throw(err error, fns ...WrapFunc) {
	if err == nil {
		return
	}
	if len(fns) > 0 {
		err = DefaultWrap(err, fns...)
	}
	if err == nil {
		return
	}
	panic(&throw{
		err: err,
	})
}
