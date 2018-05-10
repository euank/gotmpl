package gotmpl

import "errors"

var UnmatchedBraceError = errors.New("unmatched open '{'")

// UnresolvedVariableError is the error that occurs when the provided lookup
// method could not resolve a template variable
type UnresolvedVariableError struct {
	v string
}

func (u UnresolvedVariableError) Error() string {
	return "unresolved template variable: " + u.v
}
