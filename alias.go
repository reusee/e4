package e4

import (
	"errors"
	"fmt"
)

type (
	any = interface{}
)

var (
	is = errors.Is
	as = errors.As
	pt = fmt.Printf
)
