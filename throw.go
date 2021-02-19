package e4

type throw struct {
	err error
}

func (c *throw) String() string { // NOCOVER
	return c.err.Error()
}

func Throw(err error) {
	panic(&throw{
		err: err,
	})
}
