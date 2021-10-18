package e4

import "testing"

// TestingFatal returns a WrapFunc that calls t.Fatal if error occur
func TestingFatal(t *testing.T) WrapFunc {
	t.Helper()
	return WrapStacktrace.With(
		func(err error) error {
			if err == nil {
				return nil
			}
			t.Helper()
			t.Fatal(err)
			return err
		},
	)
}
