package e4

import "errors"

func Handle(errp *error, fns ...WrapFunc) {
	var err error
	// check throw error
	if p := recover(); p != nil {
		if e, ok := p.(*throw); ok {
			err = e.err
		} else {
			panic(p)
		}
	}
	// check pointed error
	if errp != nil && *errp != nil {
		if err == nil {
			// no throw error
			err = *errp
		} else {
			// wrap if not the same
			if !errors.Is(err, *errp) && !errors.Is(*errp, err) {
				err = MakeErr(err, *errp)
			}
		}
	}
	if err == nil {
		// no error
		return
	}
	// wrap
	err = Wrap(err, fns...)
	if errp != nil {
		// set pointed variable
		*errp = err
	} else {
		// throw
		Throw(err)
	}
}
