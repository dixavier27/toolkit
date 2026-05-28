package tupavendor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "internal", "tupa")
	if err := Copy(dest); err != nil {
		t.Fatal(err)
	}

	required := []string{"app.go", "recurso.go", "opcoes.go", "ganchos.go", "repositorio.go", "contexto.go", "erros.go"}
	for _, name := range required {
		p := filepath.Join(dest, name)
		info, err := os.Stat(p)
		if err != nil {
			t.Errorf("%s não criado: %v", name, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s vazio", name)
		}
	}

	// Sanity: o app.go ainda deve declarar `package tupa`.
	data, err := os.ReadFile(filepath.Join(dest, "app.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "package tupa") {
		t.Error("app.go não tem `package tupa`")
	}
}

func TestFiles(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 7 {
		t.Errorf("esperava ≥7 arquivos, got %d: %v", len(files), files)
	}
	// Nenhum deve ter sufixo .gotxt no nome reportado.
	for _, f := range files {
		if strings.HasSuffix(f, ".gotxt") {
			t.Errorf("Files() devolveu nome com .gotxt: %q", f)
		}
	}
}
