package cigen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate_Both(t *testing.T) {
	dir := t.TempDir()
	res, err := Generate(Options{Dir: dir, Name: "myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 2 {
		t.Fatalf("written = %d, want 2", len(res.Written))
	}

	for _, p := range res.Written {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("não criado: %v", err)
		}
	}

	ci, err := os.ReadFile(filepath.Join(dir, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ci), "golangci-lint-action") {
		t.Error("ci.yml não menciona golangci-lint-action")
	}
	if !strings.Contains(string(ci), "go-version: '1.26'") {
		t.Error("ci.yml não fixa go-version")
	}

	rel, _ := os.ReadFile(filepath.Join(dir, ".github", "workflows", "release.yml"))
	if !strings.Contains(string(rel), "myapp_") {
		t.Errorf("release.yml não usa name=myapp; conteúdo head=%q", string(rel[:200]))
	}
	if !strings.Contains(string(rel), "actions/upload-artifact@v4") {
		t.Error("release.yml não tem upload-artifact")
	}
	// As construções ${{ ... }} do GitHub Actions devem aparecer literalmente.
	if !strings.Contains(string(rel), "${{ matrix.goos }}") {
		t.Error("release.yml não preserva ${{ matrix.goos }}")
	}
}

func TestGenerate_OnlyCI(t *testing.T) {
	dir := t.TempDir()
	res, err := Generate(Options{Dir: dir, Only: "ci"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 1 {
		t.Errorf("written = %d, want 1", len(res.Written))
	}
	if !strings.HasSuffix(res.Written[0], "ci.yml") {
		t.Errorf("não gerou ci.yml: %v", res.Written)
	}
}

func TestGenerate_OnlyInvalid(t *testing.T) {
	dir := t.TempDir()
	_, err := Generate(Options{Dir: dir, Only: "bogus"})
	if err == nil {
		t.Fatal("esperado erro")
	}
}

func TestGenerate_NoForce(t *testing.T) {
	dir := t.TempDir()
	if _, err := Generate(Options{Dir: dir, Only: "ci"}); err != nil {
		t.Fatal(err)
	}
	res, err := Generate(Options{Dir: dir, Only: "ci"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 0 || len(res.Skipped) != 1 {
		t.Errorf("written=%d skipped=%d", len(res.Written), len(res.Skipped))
	}
}

func TestGenerate_Force(t *testing.T) {
	dir := t.TempDir()
	if _, err := Generate(Options{Dir: dir, Only: "ci"}); err != nil {
		t.Fatal(err)
	}
	res, err := Generate(Options{Dir: dir, Only: "ci", Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 1 {
		t.Errorf("written = %d, want 1 com --force", len(res.Written))
	}
}
