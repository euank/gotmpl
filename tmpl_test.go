package gotmpl

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	return len(p) - 1, nil
}

func TestMapLookupTemplate(t *testing.T) {
	s := "echo ${foo}"
	l := map[string]string{"foo": "bar"}

	res, err := TemplateString(s, MapLookup(l))
	if err != nil {
		t.Fatal(err)
	}
	if res != "echo bar" {
		t.Errorf("Expected `echo bar`, got %v", res)
	}
}

func TestEarlyEOFTemplate(t *testing.T) {
	_, err := TemplateString(`echo ${foo`, MapLookup(map[string]string{"foo": "bar"}))
	if err == nil {
		t.Error("Should fail on mismatched braces")
	}
	if err != UnmatchedBraceError {
		t.Errorf("error should be equal to UnmatchedBraceError; %v != %v", err, UnmatchedBraceError)
	}
}

func TestVarError(t *testing.T) {
	_, err := TemplateString(`echo ${foo}`, MapLookup(map[string]string{"notfoo": "bar"}))
	if err == nil {
		t.Error("Should fail on missing variable")
	}

	switch err.(type) {
	case UnresolvedVariableError:
	default:
		t.Errorf("type should be UnresolvedVariableError; was %T", err)
	}
}

func TestTemplateReturnsWriterError(t *testing.T) {
	want := errors.New("write failed")
	err := Template(strings.NewReader("content"), errorWriter{err: want}, MapLookup{})
	if !errors.Is(err, want) {
		t.Fatalf("Template() error = %v, want %v", err, want)
	}
}

func TestTemplateReturnsShortWrite(t *testing.T) {
	err := Template(strings.NewReader("content"), shortWriter{}, MapLookup{})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("Template() error = %v, want %v", err, io.ErrShortWrite)
	}
}
