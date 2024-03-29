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

// Throw checks the error and if not nil, raise a panic which will be recovered by Handle
func Throw(err error, args ...error) error {
	if err == nil {
		return nil
	}
	if len(args) > 0 {
		err = Wrap.With(args...)(err)
	}
	if err == nil {
		return nil
	}
	panic(&throw{
		err: err,
	})
}
