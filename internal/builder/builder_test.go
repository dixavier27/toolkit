package builder

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetectEntry_CmdApi(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, "cmd", "api"))
	touch(t, filepath.Join(dir, "cmd", "api", "main.go"))

	got, err := DetectEntry(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "./cmd/api" {
		t.Errorf("got %q, want ./cmd/api", got)
	}
}

func TestDetectEntry_SingleCmd(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, "cmd", "worker"))
	touch(t, filepath.Join(dir, "cmd", "worker", "main.go"))

	got, err := DetectEntry(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "./cmd/worker" {
		t.Errorf("got %q, want ./cmd/worker", got)
	}
}

func TestDetectEntry_MultipleCmds(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"api", "worker", "cli"} {
		// nota: "api" não conta como ambíguo pois é resolvido antes — mas com 3 cmds incluindo api,
		// a função retorna ./cmd/api por prioridade. Para testar ambiguidade real, omitimos "api".
		_ = name
	}
	mkdirAll(t, filepath.Join(dir, "cmd", "alpha"))
	touch(t, filepath.Join(dir, "cmd", "alpha", "main.go"))
	mkdirAll(t, filepath.Join(dir, "cmd", "beta"))
	touch(t, filepath.Join(dir, "cmd", "beta", "main.go"))

	_, err := DetectEntry(dir)
	if err == nil {
		t.Fatal("esperado erro de ambiguidade")
	}
	if !strings.Contains(err.Error(), "múltiplos cmds") {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestDetectEntry_RootMain(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "main.go"))

	got, err := DetectEntry(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "." {
		t.Errorf("got %q, want .", got)
	}
}

func TestDetectEntry_None(t *testing.T) {
	dir := t.TempDir()
	if _, err := DetectEntry(dir); err == nil {
		t.Fatal("esperado erro quando não há main.go")
	}
}

func TestBuild_Minimal(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/x\n\ngo 1.26\n")
	writeFile(t, filepath.Join(dir, "main.go"), "package main\n\nvar version = \"dev\"\n\nfunc main() { _ = version }\n")

	res, err := Build(context.Background(), Options{Dir: dir, Version: "0.0.1-test"})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if res.Out == "" {
		t.Error("Out vazio")
	}
	if res.Size <= 0 {
		t.Error("Size zero")
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(res.Out, ".exe") {
		t.Errorf("Out=%q sem .exe no windows", res.Out)
	}
	if _, err := os.Stat(res.Out); err != nil {
		t.Errorf("binário não existe: %v", err)
	}
}

func TestBuildLDFlags(t *testing.T) {
	cases := []struct {
		name string
		opts Options
		want string
	}{
		{"empty", Options{}, ""},
		{"version", Options{Version: "1.2.3"}, "-X main.version=1.2.3"},
		{"strip", Options{Strip: true}, "-s -w"},
		{"strip+version", Options{Strip: true, Version: "v1"}, "-s -w -X main.version=v1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildLDFlags(tc.opts); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// helpers
func mkdirAll(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}
func touch(t *testing.T, p string) {
	t.Helper()
	writeFile(t, p, "package main\n\nfunc main() {}\n")
}
func writeFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
