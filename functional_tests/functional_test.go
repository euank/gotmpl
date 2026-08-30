package functional_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var gotmplBinary string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "gotmpl-functional-tests-*")
	if err != nil {
		panic(err)
	}
	gotmplBinary = filepath.Join(tmpDir, "gotmpl")
	build := exec.Command("go", "build", "-o", gotmplBinary, "../cmd/gotmpl")
	if output, err := build.CombinedOutput(); err != nil {
		_, _ = os.Stderr.Write(output)
		_ = os.RemoveAll(tmpDir)
		panic(err)
	}

	code := m.Run()
	_ = os.RemoveAll(tmpDir)
	os.Exit(code)
}

func TestSimpleUsage(t *testing.T) {
	cmd := exec.Command(gotmplBinary, "./testdata/foobar.json")

	cmd.Stdin = strings.NewReader("hello ${foo}")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v, %v", string(out), err)
	}

	if string(out) != "hello bar" {
		t.Errorf("Expected `hello bar`, got %v", string(out))
	}
}

func TestSimpleFileUsage(t *testing.T) {
	t.Setenv("foo", "bar")

	cmd := exec.Command(gotmplBinary, "./testdata/template_me")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v, %v", string(out), err)
	}

	if string(out) != "hello bar" {
		t.Errorf("Expected `hello bar`, got %v", string(out))
	}
}

func TestInplacePreservesPermissions(t *testing.T) {
	t.Setenv("foo", "bar")
	file := filepath.Join(t.TempDir(), "template")
	if err := os.WriteFile(file, []byte("hello ${foo}"), 0640); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(gotmplBinary, "-inplace", file)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", output, err)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello bar" {
		t.Fatalf("content = %q, want %q", content, "hello bar")
	}
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0640 {
		t.Fatalf("permissions = %o, want 640", got)
	}
}

func TestFailedInplaceLeavesOriginal(t *testing.T) {
	file := filepath.Join(t.TempDir(), "template")
	want := "hello ${missing}"
	if err := os.WriteFile(file, []byte(want), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(gotmplBinary, "-env=false", "-inplace", file)
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("command unexpectedly succeeded: %s", output)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != want {
		t.Fatalf("content = %q, want unchanged content %q", content, want)
	}
}
