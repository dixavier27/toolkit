// Package cigen gera workflows do GitHub Actions adequados a projetos Go.
package cigen

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// Options controla a geração.
type Options struct {
	Dir       string // raiz do projeto (default: cwd)
	Only      string // "ci", "release" ou "" (ambos)
	Force     bool
	Name      string // nome do binário (default: basename do projeto)
	Entry     string // pacote a compilar (default: ./cmd/api ou .)
	GoVersion string // ex: "1.26" — default: "1.26"
}

// Result lista os arquivos gravados.
type Result struct {
	Written []string
	Skipped []string // existiam e --force=false
}

type templateData struct {
	Name      string
	Entry     string
	GoVersion string
}

// Generate grava os workflows selecionados em .github/workflows/.
func Generate(opts Options) (Result, error) {
	var res Result

	if opts.Dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return res, err
		}
		opts.Dir = cwd
	}
	absDir, err := filepath.Abs(opts.Dir)
	if err != nil {
		return res, err
	}

	if opts.Name == "" {
		opts.Name = filepath.Base(absDir)
	}
	if opts.Entry == "" {
		opts.Entry = "./..."
	}
	if opts.GoVersion == "" {
		opts.GoVersion = "1.26"
	}

	workflowsDir := filepath.Join(absDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0o755); err != nil {
		return res, err
	}

	files := selectFiles(opts.Only)
	if len(files) == 0 {
		return res, fmt.Errorf("--only inválido: %q (valores: ci, release)", opts.Only)
	}

	data := templateData{
		Name:      opts.Name,
		Entry:     opts.Entry,
		GoVersion: opts.GoVersion,
	}

	for _, name := range files {
		out := filepath.Join(workflowsDir, name)
		if _, err := os.Stat(out); err == nil && !opts.Force {
			res.Skipped = append(res.Skipped, out)
			continue
		}
		rendered, err := render("templates/"+name+".tmpl", data)
		if err != nil {
			return res, err
		}
		if err := os.WriteFile(out, rendered, 0o644); err != nil {
			return res, fmt.Errorf("write %s: %w", out, err)
		}
		res.Written = append(res.Written, out)
	}

	return res, nil
}

func selectFiles(only string) []string {
	switch strings.ToLower(strings.TrimSpace(only)) {
	case "", "all":
		return []string{"ci.yml", "release.yml"}
	case "ci":
		return []string{"ci.yml"}
	case "release":
		return []string{"release.yml"}
	default:
		return nil
	}
}

func render(name string, data templateData) ([]byte, error) {
	raw, err := templatesFS.ReadFile(name)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(filepath.Base(name)).Parse(string(raw))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
