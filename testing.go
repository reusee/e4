package e4

import "testing"

func TestingFatal(t *testing.T) WrapFunc {
	t.Helper()
	return func(err error) error {
		if err == nil {
			return nil
		}
		t.Helper()
		t.Fatal(err)
		return err
	}
}
