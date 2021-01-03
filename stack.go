package e4

func Stack(fn func() error) WrapFunc {
	return func(err error) error {
		if e := fn(); e != nil {
			return Error{
				Err:  e,
				Prev: err,
			}
		}
		return err
	}
}
