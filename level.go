package e4

type Level int

const (
	Critical Level = iota + 1
	Informational
	Debug
)

const ErrorLevel = Informational

type ErrorLeveler interface {
	ErrorLevel() Level
}
