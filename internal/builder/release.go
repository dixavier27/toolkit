package builder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Target descreve um par GOOS/GOARCH para cross-compile.
type Target struct {
	GOOS, GOARCH string
}

// DefaultTargets retorna a matriz padrão de releases multi-OS.
func DefaultTargets() []Target {
	return []Target{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}
}

// ReleaseOptions controla um build multi-target.
type ReleaseOptions struct {
	Dir       string // raiz do projeto (default: cwd)
	Entry     string // entrypoint (default: detectado)
	OutDir    string // diretório de saída (default: "release/")
	Name      string // nome base do binário (default: basename do projeto)
	Version   string // versão injetada (default: git describe)
	Strip     bool   // -ldflags '-s -w'
	Verbose   bool
	Targets   []Target // se vazio: DefaultTargets()
	KeepGoing bool     // continua em caso de falha individual
}

// ReleaseResult agrega os artefatos gerados.
type ReleaseResult struct {
	Artifacts     []Result
	ChecksumsFile string
	Failed        []TargetError
}

// TargetError descreve uma falha de target individual.
type TargetError struct {
	Target Target
	Err    error
}

// Release compila o projeto para todos os targets configurados e gera um
// arquivo checksums.txt no OutDir.
func Release(ctx context.Context, opts ReleaseOptions) (ReleaseResult, error) {
	var res ReleaseResult

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
	opts.Dir = absDir

	if opts.OutDir == "" {
		opts.OutDir = "release"
	}
	if !filepath.IsAbs(opts.OutDir) {
		opts.OutDir = filepath.Join(absDir, opts.OutDir)
	}
	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return res, err
	}

	if opts.Name == "" {
		opts.Name = filepath.Base(absDir)
	}
	if opts.Version == "" {
		opts.Version = GitVersion(absDir)
	}
	if len(opts.Targets) == 0 {
		opts.Targets = DefaultTargets()
	}

	for _, t := range opts.Targets {
		out := filepath.Join(opts.OutDir, artifactName(opts.Name, t, opts.Version))
		r, err := Build(ctx, Options{
			Dir:     opts.Dir,
			Entry:   opts.Entry,
			Out:     out,
			GOOS:    t.GOOS,
			GOARCH:  t.GOARCH,
			Version: opts.Version,
			Strip:   opts.Strip,
			Verbose: opts.Verbose,
		})
		if err != nil {
			res.Failed = append(res.Failed, TargetError{Target: t, Err: err})
			if opts.KeepGoing {
				continue
			}
			return res, fmt.Errorf("build %s/%s: %w", t.GOOS, t.GOARCH, err)
		}
		res.Artifacts = append(res.Artifacts, r)
	}

	if len(res.Artifacts) == 0 {
		return res, fmt.Errorf("nenhum artefato gerado")
	}

	checksums, err := writeChecksums(opts.OutDir, res.Artifacts)
	if err != nil {
		return res, fmt.Errorf("checksums: %w", err)
	}
	res.ChecksumsFile = checksums

	return res, nil
}

// artifactName produz nomes consistentes: <name>_<version>_<os>_<arch>[.exe]
// Se version == "dev" ou vazio, omite o segmento de versão.
func artifactName(name string, t Target, version string) string {
	suffix := ""
	if t.GOOS == "windows" {
		suffix = ".exe"
	}
	if version == "" || version == "dev" {
		return fmt.Sprintf("%s_%s_%s%s", name, t.GOOS, t.GOARCH, suffix)
	}
	return fmt.Sprintf("%s_%s_%s_%s%s", name, version, t.GOOS, t.GOARCH, suffix)
}

func writeChecksums(outDir string, artifacts []Result) (string, error) {
	type entry struct {
		hash string
		file string
	}
	var entries []entry
	for _, a := range artifacts {
		h, err := sha256File(a.Out)
		if err != nil {
			return "", err
		}
		entries = append(entries, entry{hash: h, file: filepath.Base(a.Out)})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].file < entries[j].file })

	path := filepath.Join(outDir, "checksums.txt")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	for _, e := range entries {
		if _, err := fmt.Fprintf(f, "%s  %s\n", e.hash, e.file); err != nil {
			return "", err
		}
	}
	return path, nil
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
