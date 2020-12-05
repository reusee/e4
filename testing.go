package e4

import "testing"

func TestingFatal(t *testing.T) WrapFunc {
	t.Helper()
	return func(err error) error {
		t.Helper()
		t.Fatal(err)
		return err
	}
}
