// Package builder encapsula a invocação de `go build` com opções comuns
// (injeção de versão, strip, cross-compile). Reutilizado por `eco build`
// e `eco release`.
package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Options controla uma única invocação de `go build`.
type Options struct {
	// Dir é o diretório raiz do projeto (onde está go.mod). Default: cwd.
	Dir string

	// Entry é o pacote a compilar (ex.: "./cmd/api"). Se vazio, é detectado.
	Entry string

	// Out é o caminho do binário gerado. Se vazio, default: bin/<nome>[.exe].
	Out string

	// GOOS / GOARCH para cross-compile. Vazio = nativo.
	GOOS, GOARCH string

	// Version é injetada em `main.version` via -ldflags.
	// Se vazio, builder tenta `git describe`.
	Version string

	// Strip aplica `-s -w` ao linker (binário menor).
	Strip bool

	// Verbose ecoa o comando executado.
	Verbose bool
}

// Result descreve um build bem-sucedido.
type Result struct {
	Out    string
	GOOS   string
	GOARCH string
	Size   int64 // bytes
}

// Build compila o projeto em opts.Dir conforme as opções. Devolve o caminho
// do binário gerado (absoluto).
func Build(ctx context.Context, opts Options) (Result, error) {
	if opts.Dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return Result{}, err
		}
		opts.Dir = cwd
	}
	absDir, err := filepath.Abs(opts.Dir)
	if err != nil {
		return Result{}, err
	}

	if opts.Entry == "" {
		entry, err := DetectEntry(absDir)
		if err != nil {
			return Result{}, err
		}
		opts.Entry = entry
	}

	if opts.Version == "" {
		opts.Version = GitVersion(absDir)
	}

	if opts.Out == "" {
		opts.Out = defaultOut(absDir, opts.GOOS)
	} else if !filepath.IsAbs(opts.Out) {
		opts.Out = filepath.Join(absDir, opts.Out)
	}
	if err := os.MkdirAll(filepath.Dir(opts.Out), 0o755); err != nil {
		return Result{}, err
	}

	args := []string{"build"}
	if ldflags := buildLDFlags(opts); ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, "-o", opts.Out, opts.Entry)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = absDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if opts.GOOS != "" {
		cmd.Env = append(cmd.Env, "GOOS="+opts.GOOS)
	}
	if opts.GOARCH != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+opts.GOARCH)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "→ %s %s\n", "go", strings.Join(args, " "))
	}

	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("go build: %w", err)
	}

	info, err := os.Stat(opts.Out)
	if err != nil {
		return Result{}, fmt.Errorf("stat output: %w", err)
	}

	res := Result{
		Out:    opts.Out,
		GOOS:   resolveOS(opts.GOOS),
		GOARCH: resolveArch(opts.GOARCH),
		Size:   info.Size(),
	}
	return res, nil
}

// DetectEntry escolhe o entrypoint padrão para o projeto.
// Prioridade: cmd/api/main.go → cmd/<nome>/main.go (único) → main.go raiz.
func DetectEntry(dir string) (string, error) {
	api := filepath.Join(dir, "cmd", "api", "main.go")
	if fileExists(api) {
		return "./cmd/api", nil
	}

	cmdDir := filepath.Join(dir, "cmd")
	if entries, err := os.ReadDir(cmdDir); err == nil {
		var subdirs []string
		for _, e := range entries {
			if e.IsDir() && fileExists(filepath.Join(cmdDir, e.Name(), "main.go")) {
				subdirs = append(subdirs, e.Name())
			}
		}
		if len(subdirs) == 1 {
			return "./cmd/" + subdirs[0], nil
		}
		if len(subdirs) > 1 {
			return "", fmt.Errorf("múltiplos cmds em ./cmd/ (%s) — passe --entry", strings.Join(subdirs, ", "))
		}
	}

	if fileExists(filepath.Join(dir, "main.go")) {
		return ".", nil
	}
	return "", fmt.Errorf("nenhum main.go encontrado (procurei em ./cmd/api/, ./cmd/*/, .)")
}

// GitVersion devolve `git describe --tags --always --dirty`. Se git não estiver
// disponível ou não houver repo, devolve "dev".
func GitVersion(dir string) string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "dev"
	}
	v := strings.TrimSpace(string(out))
	if v == "" {
		return "dev"
	}
	return v
}

func buildLDFlags(opts Options) string {
	var parts []string
	if opts.Strip {
		parts = append(parts, "-s", "-w")
	}
	if opts.Version != "" {
		parts = append(parts, fmt.Sprintf("-X main.version=%s", opts.Version))
	}
	return strings.Join(parts, " ")
}

func defaultOut(dir, goos string) string {
	name := filepath.Base(dir)
	if resolveOS(goos) == "windows" {
		name += ".exe"
	}
	return filepath.Join(dir, "bin", name)
}

func resolveOS(goos string) string {
	if goos != "" {
		return goos
	}
	return runtime.GOOS
}

func resolveArch(goarch string) string {
	if goarch != "" {
		return goarch
	}
	return runtime.GOARCH
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
