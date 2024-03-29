package e4

import (
	"io"
	"testing"
)

func TestInfo(t *testing.T) {
	TestWrapFunc(t, Info("foo"))

	info := Info("foo %s", "bar")(io.EOF)
	if info.Error() != "foo bar\nEOF" {
		t.Fatalf("got %s", info.Error())
	}
	if !is(info, io.EOF) {
		t.Fatal()
	}
}
