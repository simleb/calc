package calc

import (
	"fmt"
)

// A calcError is a pretty printing error for calc.
type calcError struct {
	err string
	exp string
	loc int
}

// Error returns the formatted error.
func (e calcError) Error() string {
	return fmt.Sprintf("calc: %s\n      %s\n%*s", e.err, e.exp, e.loc+7, "^")
}

// makeError creates a calcError.
func makeError(exp string, loc int, desc string, args ...interface{}) calcError {
	return calcError{
		err: fmt.Sprintf(desc, args...),
		exp: exp,
		loc: loc,
	}
}
