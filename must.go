package e4

func Must(err error, fns ...WrapFunc) error {
	if err == nil {
		return nil
	}
	err = DefaultWrap(err, fns...)
	panic(err)
}
