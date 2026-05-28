package builder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArtifactName(t *testing.T) {
	cases := []struct {
		name, version string
		target        Target
		want          string
	}{
		{"eco", "v0.1.0", Target{"linux", "amd64"}, "eco_v0.1.0_linux_amd64"},
		{"eco", "v0.1.0", Target{"windows", "amd64"}, "eco_v0.1.0_windows_amd64.exe"},
		{"eco", "", Target{"darwin", "arm64"}, "eco_darwin_arm64"},
		{"eco", "dev", Target{"linux", "arm64"}, "eco_linux_arm64"},
	}
	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			if got := artifactName(tc.name, tc.target, tc.version); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDefaultTargets(t *testing.T) {
	ts := DefaultTargets()
	if len(ts) == 0 {
		t.Fatal("empty")
	}
	// sanity: contém windows/amd64
	found := false
	for _, t2 := range ts {
		if t2.GOOS == "windows" && t2.GOARCH == "amd64" {
			found = true
			break
		}
	}
	if !found {
		t.Error("windows/amd64 ausente da matriz default")
	}
}

func TestRelease_Subset(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/r\n\ngo 1.26\n")
	writeFile(t, filepath.Join(dir, "main.go"), "package main\n\nvar version = \"dev\"\n\nfunc main() { _ = version }\n")

	// Apenas linux/amd64 + windows/amd64 para manter o teste rápido.
	res, err := Release(context.Background(), ReleaseOptions{
		Dir:     dir,
		Name:    "r",
		Version: "v0.0.1",
		Targets: []Target{
			{"linux", "amd64"},
			{"windows", "amd64"},
		},
	})
	if err != nil {
		t.Fatalf("release: %v", err)
	}
	if len(res.Artifacts) != 2 {
		t.Errorf("artifacts = %d, want 2", len(res.Artifacts))
	}
	if res.ChecksumsFile == "" {
		t.Fatal("checksums vazio")
	}
	data, err := os.ReadFile(res.ChecksumsFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), "\n"); got != 2 {
		t.Errorf("checksums tem %d linhas, want 2; conteúdo=%q", got, string(data))
	}
}
