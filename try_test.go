package e4

import (
	"io"
	"testing"
)

func TestTry(t *testing.T) {
	fn := func() (int, error) {
		return 42, io.EOF
	}
	var err error
	func() {
		defer Handle(&err)
		Try(fn())(Check)
	}()
	if err != io.EOF {
		t.Fatal()
	}

	fn = func() (int, error) {
		return 42, nil
	}
	err = nil
	var i int
	func() {
		defer Handle(&err)
		i = Try(fn())(Check)
	}()
	if err != nil {
		t.Fatal()
	}
	if i != 42 {
		t.Fatal()
	}
}

func TestTry2(t *testing.T) {
	fn := func() (int, string, error) {
		return 42, "foo", io.EOF
	}
	var err error
	func() {
		defer Handle(&err)
		Try2(fn())(Check)
	}()
	if err != io.EOF {
		t.Fatal()
	}

	fn = func() (int, string, error) {
		return 42, "foo", nil
	}
	err = nil
	var i int
	var s string
	func() {
		defer Handle(&err)
		i, s = Try2(fn())(Check)
	}()
	if err != nil {
		t.Fatal()
	}
	if i != 42 {
		t.Fatal()
	}
	if s != "foo" {
		t.Fatal()
	}
}

func TestTry3(t *testing.T) {
	fn := func() (int, string, int8, error) {
		return 42, "foo", 1, io.EOF
	}
	var err error
	func() {
		defer Handle(&err)
		Try3(fn())(Check)
	}()
	if err != io.EOF {
		t.Fatal()
	}

	fn = func() (int, string, int8, error) {
		return 42, "foo", 1, nil
	}
	err = nil
	var i int
	var s string
	var i8 int8
	func() {
		defer Handle(&err)
		i, s, i8 = Try3(fn())(Check)
	}()
	if err != nil {
		t.Fatal()
	}
	if i != 42 {
		t.Fatal()
	}
	if s != "foo" {
		t.Fatal()
	}
	if i8 != 1 {
		t.Fatal()
	}
}

func TestTry4(t *testing.T) {
	fn := func() (int, string, int8, int16, error) {
		return 42, "foo", 1, 2, io.EOF
	}
	var err error
	func() {
		defer Handle(&err)
		Try4(fn())(Check)
	}()
	if err != io.EOF {
		t.Fatal()
	}

	fn = func() (int, string, int8, int16, error) {
		return 42, "foo", 1, 2, nil
	}
	err = nil
	var i int
	var s string
	var i8 int8
	var i16 int16
	func() {
		defer Handle(&err)
		i, s, i8, i16 = Try4(fn())(Check)
	}()
	if err != nil {
		t.Fatal()
	}
	if i != 42 {
		t.Fatal()
	}
	if s != "foo" {
		t.Fatal()
	}
	if i8 != 1 {
		t.Fatal()
	}
	if i16 != 2 {
		t.Fatal()
	}
}
