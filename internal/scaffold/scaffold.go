// Package scaffold materializa a estrutura padrão de um projeto de API REST em Go
// a partir de templates embarcados.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates
var templatesFS embed.FS

// Options controla a geração do projeto.
type Options struct {
	Name   string // nome do diretório destino
	Module string // go module path
	Force  bool   // sobrescrever diretório não vazio
}

// Data é o contexto passado para cada template.
type Data struct {
	Name   string
	Module string
}

// Generate cria o projeto em ./<opts.Name>. Devolve o caminho absoluto criado.
func Generate(opts Options) (string, error) {
	if opts.Name == "" {
		return "", fmt.Errorf("nome do projeto vazio")
	}
	if opts.Module == "" {
		opts.Module = opts.Name
	}

	dest, err := filepath.Abs(opts.Name)
	if err != nil {
		return "", err
	}

	if err := ensureDestDir(dest, opts.Force); err != nil {
		return "", err
	}

	data := Data{Name: opts.Name, Module: opts.Module}

	walkErr := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "templates" {
			return nil
		}

		rel := strings.TrimPrefix(path, "templates/")
		// Especial: gitignore.tmpl → .gitignore
		outRel := strings.TrimSuffix(rel, ".tmpl")
		if outRel == "gitignore" {
			outRel = ".gitignore"
		}
		outPath := filepath.Join(dest, outRel)

		if d.IsDir() {
			return os.MkdirAll(outPath, 0o755)
		}

		return renderFile(path, outPath, data)
	})
	if walkErr != nil {
		return "", walkErr
	}

	if err := runGoModTidy(dest); err != nil {
		// Não fatal: avisar mas seguir adiante.
		fmt.Fprintf(os.Stderr, "aviso: `go mod tidy` falhou em %s: %v\n", dest, err)
	}

	return dest, nil
}

func ensureDestDir(dest string, force bool) error {
	info, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return os.MkdirAll(dest, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s existe e não é um diretório", dest)
	}
	// Diretório existe: ok se vazio, ou se --force.
	entries, err := os.ReadDir(dest)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !force {
		return fmt.Errorf("diretório %s não está vazio (use --force para sobrescrever)", dest)
	}
	return nil
}

func renderFile(srcPath, outPath string, data Data) error {
	raw, err := templatesFS.ReadFile(srcPath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(filepath.Base(srcPath)).Parse(string(raw))
	if err != nil {
		return fmt.Errorf("parse %s: %w", srcPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func runGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
