package e4

// Must checks the error, if not nil, raise a panic which will not be catched by Handle
func Must(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = Wrap(fns...)(err)
	panic(err)
}

var Fatal = Must
