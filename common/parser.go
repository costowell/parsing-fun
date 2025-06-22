package common

type Parser interface {
	Parse(input string) bool
}
