// Package tupavendor embarca uma cópia do source do tupa-go para que
// projetos gerados pelo eco com --tupa funcionem sem dependência externa.
//
// Os arquivos vivem em source/ com extensão .gotxt para não serem
// compilados como parte do pacote tupavendor (eles declaram package tupa).
// Copy() materializa cada source/<nome>.gotxt como <dest>/<nome>.go.
package tupavendor

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:source
var sourceFS embed.FS

// Copy materializa os arquivos vendored em destDir. Os .gotxt viram .go;
// LICENSE é copiado como está.
func Copy(destDir string) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	return fs.WalkDir(sourceFS, "source", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "source" {
			return nil
		}
		name := strings.TrimPrefix(path, "source/")
		outName := strings.TrimSuffix(name, ".gotxt")
		if outName != "LICENSE" && !strings.HasSuffix(outName, ".go") {
			outName += ".go"
		}
		if name == "LICENSE" {
			outName = "LICENSE"
		}
		out := filepath.Join(destDir, outName)

		if d.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		data, err := sourceFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", out, err)
		}
		return nil
	})
}

// Files retorna a lista de arquivos que serão materializados (relativo a destDir).
// Útil para mensagens no CLI ("criado: internal/tupa/recurso.go").
func Files() ([]string, error) {
	var out []string
	err := fs.WalkDir(sourceFS, "source", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := strings.TrimPrefix(path, "source/")
		out = append(out, strings.TrimSuffix(name, ".gotxt"))
		return nil
	})
	return out, err
}
