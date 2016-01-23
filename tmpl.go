package gotmpl

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

func Template(r io.Reader, w io.Writer, lookup Lookup) error {
	bufReader := bufio.NewReader(r)
	inTemplate := false
	varName := ""
	for {
		b, err := bufReader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Here there be dragons; TODO
			continue
		}

		if inTemplate {
			if b == '}' {
				inTemplate = false
				val, ok := lookup.Resolve(varName)
				if !ok {
					return errors.New("Could not resolve variable: " + varName)
				}
				w.Write([]byte(val))
			} else {
				varName += string(b)
			}
			continue
		}

		if b == '\\' {
			nb, err := bufReader.Peek(1)
			if err == io.EOF {
				w.Write([]byte{b})
				break
			}
			if nb[0] == byte('\\') {
				// \\ escape
				w.Write([]byte{b})
				bufReader.ReadByte()
				continue
			}
			if nb[0] == '$' {
				// \$ escape
				w.Write([]byte("$"))
				bufReader.ReadByte()
				continue
			}
		}

		if b == '$' {
			nb, err := bufReader.Peek(1)
			if err == io.EOF {
				w.Write([]byte{b})
				break
			}
			if nb[0] == '{' {
				inTemplate = true
				varName = ""
				bufReader.ReadByte()
				continue
			}
		}

		w.Write([]byte{b})
	}
	return nil
}

func TemplateString(templateString string, lookup Lookup) (string, error) {
	var out bytes.Buffer
	err := Template(strings.NewReader(templateString), &out, lookup)
	return out.String(), err
}
