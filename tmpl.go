// Package gotmpl provides a simple library for stupid-simple templating.
// This templating is limited to variable substitution only. The only special characters are `\` and `$`.
// Each can be escaped with a backslash as `\\` and `\$` respectively.
// A variable reference takes the form `${variable}`. Variable names may contain
// any character except `}`.
// Please see the examples and README for more details on usage of this library.
package gotmpl

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Template parses any variables from a given reader and outputs processed data
// to the given writer. Variables are replaced based on the result of the
// provided lookup.
func Template(r io.Reader, w io.Writer, lookup Lookup) error {
	bufReader := bufio.NewReader(r)
	write := func(p []byte) error {
		n, err := w.Write(p)
		if err != nil {
			return err
		}
		if n != len(p) {
			return io.ErrShortWrite
		}
		return nil
	}
	inTemplate := false
	varName := ""
	for {
		b, err := bufReader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if inTemplate {
			if b == '}' {
				inTemplate = false
				val, ok := lookup.Resolve(varName)
				if !ok {
					return UnresolvedVariableError{v: varName}
				}
				if err := write([]byte(val)); err != nil {
					return err
				}
			} else {
				varName += string(b)
			}
			continue
		}

		if b == '\\' {
			nb, err := bufReader.Peek(1)
			if err == io.EOF {
				if err := write([]byte{b}); err != nil {
					return err
				}
				break
			} else if err != nil {
				return err
			}
			if nb[0] == byte('\\') {
				// \\ escape
				if err := write([]byte{b}); err != nil {
					return err
				}
				bufReader.ReadByte()
				continue
			}
			if nb[0] == '$' {
				// \$ escape
				if err := write([]byte("$")); err != nil {
					return err
				}
				bufReader.ReadByte()
				continue
			}
		}

		if b == '$' {
			nb, err := bufReader.Peek(1)
			if err == io.EOF {
				if err := write([]byte{b}); err != nil {
					return err
				}
				break
			} else if err != nil {
				return err
			}
			if nb[0] == '{' {
				inTemplate = true
				varName = ""
				bufReader.ReadByte()
				continue
			}
		}

		if err := write([]byte{b}); err != nil {
			return err
		}
	}
	if inTemplate {
		return UnmatchedBraceError
	}
	return nil
}

// TemplateString is a convenience function to template a given input string in memory
func TemplateString(templateString string, lookup Lookup) (string, error) {
	var out bytes.Buffer
	err := Template(strings.NewReader(templateString), &out, lookup)
	return out.String(), err
}
