package e4

import "testing"

func TestTesting(t *testing.T) {
	testWrapFunc(t, TestingFatal(t))
}
