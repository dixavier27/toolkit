// Package gocheck detecta a toolchain Go local e gera instruções de instalação
// adequadas ao sistema operacional do usuário.
package gocheck

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// MinMajor / MinMinor é a versão mínima exigida do Go (1.26).
const (
	MinMajor = 1
	MinMinor = 26
)

// Result resume o resultado da verificação do Go.
type Result struct {
	Found     bool
	Version   string // ex: "1.23.4"
	Path      string // caminho do binário go
	MeetsMin  bool
	RawOutput string
}

// Check executa `go version` e analisa a saída.
func Check() Result {
	r := Result{}

	path, err := exec.LookPath("go")
	if err != nil {
		return r
	}
	r.Path = path

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return r
	}
	r.RawOutput = strings.TrimSpace(string(out))
	r.Found = true

	// Exemplo de saída: "go version go1.23.4 linux/amd64"
	re := regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)
	m := re.FindStringSubmatch(r.RawOutput)
	if len(m) >= 3 {
		major, _ := strconv.Atoi(m[1])
		minor, _ := strconv.Atoi(m[2])
		r.Version = strings.TrimPrefix(m[0], "go")
		r.MeetsMin = major > MinMajor || (major == MinMajor && minor >= MinMinor)
	}

	return r
}

// InstallInstructions devolve um texto multilinha com instruções de instalação
// por sistema operacional.
func InstallInstructions() string {
	switch runtime.GOOS {
	case "windows":
		return strings.Join([]string{
			"Go não encontrado no PATH.",
			"",
			"Instalação no Windows:",
			"  • winget install --id GoLang.Go",
			"  • ou baixe o instalador em https://go.dev/dl/",
			"",
			"Após instalar, abra um novo terminal e rode `eco doctor` novamente.",
		}, "\n")
	case "darwin":
		return strings.Join([]string{
			"Go não encontrado no PATH.",
			"",
			"Instalação no macOS:",
			"  • brew install go",
			"  • ou baixe o pacote em https://go.dev/dl/",
			"",
			"Após instalar, rode `eco doctor` novamente.",
		}, "\n")
	case "linux":
		return strings.Join([]string{
			"Go não encontrado no PATH.",
			"",
			"Instalação no Linux:",
			"  • Pacote oficial (recomendado): https://go.dev/dl/",
			"  • Debian/Ubuntu: sudo apt install golang-go  (versão pode estar desatualizada)",
			"  • Fedora:        sudo dnf install golang",
			"  • Arch:          sudo pacman -S go",
			"",
			"Após instalar, rode `eco doctor` novamente.",
		}, "\n")
	default:
		return fmt.Sprintf("Go não encontrado. Veja https://go.dev/dl/ para instruções de instalação em %s.", runtime.GOOS)
	}
}
