package e4

import "io"

func WithClose(c io.Closer) WrapFunc {
	return func(prev error) error {
		if err := c.Close(); err != nil {
			return With(err)(prev)
		}
		return prev
	}
}
