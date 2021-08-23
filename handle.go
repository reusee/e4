package e4

import "errors"

// Handle is for error handling
//
// Error raised by Throw will be catched if any.
// If errp point to non-nil error, the error will be chained.
// If the result error is not nil, wrap functions will be applied.
// The result error will be assigned to errp if errp is not nil, otherwise Throw will be raised.
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
	err = Wrap(fns...)(err)
	if errp != nil {
		// set pointed variable
		*errp = err
	} else {
		// throw
		Throw(err)
	}
}
