package e4

import "testing"

func TestTesting(t *testing.T) {
	TestWrapFunc(t, TestingFatal(t))
}
