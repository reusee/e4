package e4

func Try[R any](r R, err error) func(fn CheckFunc) R {
	return func(fn CheckFunc) R {
		fn(err)
		return r
	}
}

func Try2[R, R2 any](r R, r2 R2, err error) func(fn CheckFunc) (R, R2) {
	return func(fn CheckFunc) (R, R2) {
		fn(err)
		return r, r2
	}
}

func Try3[R, R2, R3 any](r R, r2 R2, r3 R3, err error) func(fn CheckFunc) (R, R2, R3) {
	return func(fn CheckFunc) (R, R2, R3) {
		fn(err)
		return r, r2, r3
	}
}

func Try4[R, R2, R3, R4 any](r R, r2 R2, r3 R3, r4 R4, err error) func(fn CheckFunc) (R, R2, R3, R4) {
	return func(fn CheckFunc) (R, R2, R3, R4) {
		fn(err)
		return r, r2, r3, r4
	}
}
