package gocheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGoMod_Missing(t *testing.T) {
	dir := t.TempDir()
	c := CheckGoMod(dir)
	if c.Status != StatusMissing {
		t.Fatalf("status = %v, want StatusMissing", c.Status)
	}
}

func TestCheckGoMod_Valid(t *testing.T) {
	dir := t.TempDir()
	content := "module example.com/demo\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	c := CheckGoMod(dir)
	if c.Status != StatusOK {
		t.Fatalf("status = %v, want StatusOK; message=%q", c.Status, c.Message)
	}
	if c.Message != "módulo example.com/demo" {
		t.Errorf("message = %q", c.Message)
	}
}

func TestCheckGoMod_NoModuleDirective(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("go 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := CheckGoMod(dir)
	if c.Status != StatusWarn {
		t.Fatalf("status = %v, want StatusWarn", c.Status)
	}
}

func TestExtractModule(t *testing.T) {
	cases := map[string]string{
		"module foo/bar\ngo 1.22":               "foo/bar",
		"go 1.22\nmodule example.com/x":         "example.com/x",
		"  module   spaces/in/front  \ngo 1.22": "spaces/in/front",
		"go 1.22":                               "",
	}
	for input, want := range cases {
		if got := extractModule(input); got != want {
			t.Errorf("extractModule(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCheckGo_Smoke(t *testing.T) {
	// Não asserta versão específica; só que CheckGo não panica e devolve algum status.
	c := CheckGo()
	if c.Name != "go" {
		t.Errorf("name = %q", c.Name)
	}
	if c.Status == StatusOK && c.Version == "" {
		t.Errorf("status OK mas version vazia: %+v", c)
	}
}
