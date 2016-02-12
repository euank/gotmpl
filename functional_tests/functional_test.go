package functional_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSimpleUsage(t *testing.T) {
	cmd := exec.Command("../bin/gotmpl", "./testdata/foobar.json")

	cmd.Stdin = strings.NewReader("hello ${foo}")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v, %v", string(out), err)
	}

	if string(out) != "hello bar" {
		t.Errorf("Expected `hello bar`, got %v", string(out))
	}
}
