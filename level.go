package e4

type Level int

const (
	CriticalLevel Level = iota + 1
	InfoLevel
	DebugLevel
)

var ErrorLevel = InfoLevel

type ErrorLeveler interface {
	ErrorLevel() Level
}
